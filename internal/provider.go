package internal

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
	"golang.org/x/oauth2"
)

func New() provider.Provider {
	return &woodpeckerProvider{}
}

type woodpeckerProvider struct {
	config providerConfig
	client woodpecker.Client
	self   *woodpecker.User
}

func (p *woodpeckerProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "woodpecker"
}

func (p *woodpeckerProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"server": {
				Optional: true,
				Type:     types.StringType,
			},
			"token": {
				Optional: true,
				Type:     types.StringType,
			},
			"verify": {
				Optional: true,
				Type:     types.BoolType,
			},
		},
	}, nil
}

func (p *woodpeckerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceRepository,
		NewDataSourceRepositoryCron,
		NewDataSourceRepositorySecret,
		NewDataSourceSelf,
	}
}

func (p *woodpeckerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRepositoryResource,
		NewRepositoryCronResource,
		NewRepositorySecretResource,
	}
}

func (p *woodpeckerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.config = p.createProviderConfiguration(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	p.client, p.self = p.createClient(ctx, p.config, resp)

	resp.DataSourceData = p
	resp.ResourceData = p
}

type providerConfig struct {
	Server types.String `tfsdk:"server"`
	Token  types.String `tfsdk:"token"`
	Verify types.Bool   `tfsdk:"verify"`
}

func (p *woodpeckerProvider) createProviderConfiguration(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) providerConfig {
	var config providerConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return config
	}

	if config.Server.IsNull() {
		config.Server = types.String{Value: os.Getenv("WOODPECKER_SERVER")}
	}

	if config.Token.IsNull() {
		config.Token = types.String{Value: os.Getenv("WOODPECKER_TOKEN")}
	}

	if config.Verify.IsNull() {
		config.Verify = types.Bool{Value: os.Getenv("WOODPECKER_VERIFY") != "0"}
	}

	return config
}

func (p *woodpeckerProvider) createClient(
	ctx context.Context,
	config providerConfig,
	resp *provider.ConfigureResponse,
) (woodpecker.Client, *woodpecker.User) {

	oauth_config := new(oauth2.Config)

	authenticator := oauth_config.Client(ctx, &oauth2.Token{
		AccessToken: config.Token.ValueString(),
	})

	client := woodpecker.NewClient(config.Server.ValueString(), authenticator)

	self, err := client.Self()

	if err != nil {
		resp.Diagnostics.AddError("Unable to login", err.Error())
		return nil, nil
	}

	return client, self
}
