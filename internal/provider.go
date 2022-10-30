package internal

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
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

func (p *woodpeckerProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"woodpecker_repository": dataSourceRepositoryType{},
		"woodpecker_self":       dataSourceSelfType{},
	}, nil
}

func (p *woodpeckerProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"woodpecker_repository": resourceRepositoryType{},
	}, nil
}

func (p *woodpeckerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	p.config = p.createProviderConfiguration(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	p.client, p.self = p.createClient(ctx, p.config, resp)
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

	if config.Server.Null {
		config.Server = types.String{Value: os.Getenv("WOODPECKER_SERVER")}
	}

	if config.Token.Null {
		config.Token = types.String{Value: os.Getenv("WOODPECKER_TOKEN")}
	}

	if config.Verify.Null {
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
		AccessToken: config.Token.Value,
	})

	client := woodpecker.NewClient(config.Server.Value, authenticator)

	self, err := client.Self()

	if err != nil {
		resp.Diagnostics.AddError("Unable to login", err.Error())
		return nil, nil
	}

	return client, self
}
