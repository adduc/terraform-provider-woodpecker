package internal

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/woodpecker-ci/woodpecker/woodpecker-go/woodpecker"
	"golang.org/x/oauth2"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	config providerConfig
	client woodpecker.Client
	self   *woodpecker.User
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"woodpecker_repo": dataSourceRepoType{},
	}, nil
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {

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

func (p *provider) createProviderConfiguration(
	ctx context.Context,
	req tfsdk.ConfigureProviderRequest,
	resp *tfsdk.ConfigureProviderResponse,
) providerConfig {
	var config providerConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return config
	}

	if config.Server.Null {
		config.Server = types.String{
			Value: os.Getenv("WOODPECKER_SERVER"),
		}
	}

	if config.Token.Null {
		config.Token = types.String{
			Value: os.Getenv("WOODPECKER_TOKEN"),
		}
	}

	if config.Verify.Null {
		config.Verify = types.Bool{
			Value: os.Getenv("WOODPECKER_VERIFY") == "1",
		}
	}

	return config
}

func (p *provider) createClient(
	ctx context.Context,
	config providerConfig,
	resp *tfsdk.ConfigureProviderResponse,
) (woodpecker.Client, *woodpecker.User) {

	oauth_config := new(oauth2.Config)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.Verify.Value,
	}

	authenticator := oauth_config.Client(ctx, &oauth2.Token{
		AccessToken: config.Token.Value,
	})

	trans, _ := authenticator.Transport.(*oauth2.Transport)
	trans.Base = &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
	}

	client := woodpecker.NewClient(config.Server.Value, authenticator)

	self, err := client.Self()

	if err != nil {
		resp.Diagnostics.AddError("Unable to login", err.Error())
		return nil, nil
	}

	return client, self
}
