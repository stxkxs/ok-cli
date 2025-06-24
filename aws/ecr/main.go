package ecr

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/stxkxs/ok-cli/logger"
	"github.com/stxkxs/ok-cli/terminal"
	"os"
	"strings"
	"time"
)

type PrivateRepository struct {
	Name       string            `mapstructure:"name"`
	Version    string            `mapstructure:"version"`
	ScanOnPush bool              `mapstructure:"scanOnPush"`
	Mutability string            `mapstructure:"mutability"`
	Tags       map[string]string `mapstructure:"tags"`
}

type PublicRepository struct {
	Name             string            `mapstructure:"name"`
	Version          string            `mapstructure:"version"`
	ScanOnPush       bool              `mapstructure:"scanOnPush"`
	Mutability       string            `mapstructure:"mutability"`
	Alias            string            `mapstructure:"alias"`
	Description      string            `mapstructure:"description"`
	About            string            `mapstructure:"about"`
	Usage            string            `mapstructure:"usage"`
	Architectures    []string          `mapstructure:"architectures"`
	OperatingSystems []string          `mapstructure:"operatingSystems"`
	Tags             map[string]string `mapstructure:"tags"`
}

type PrivateDockerImage struct {
	Name       string            `mapstructure:"name"`
	Version    string            `mapstructure:"version"`
	ScanOnPush bool              `mapstructure:"scanOnPush"`
	Mutability string            `mapstructure:"mutability"`
	Dockerfile string            `mapstructure:"dockerfile"`
	Context    string            `mapstructure:"context"`
	Tags       map[string]string `mapstructure:"tags"`
}

type PublicDockerImage struct {
	Name             string            `mapstructure:"name"`
	Version          string            `mapstructure:"version"`
	ScanOnPush       bool              `mapstructure:"scanOnPush"`
	Mutability       string            `mapstructure:"mutability"`
	Dockerfile       string            `mapstructure:"dockerfile"`
	Context          string            `mapstructure:"context"`
	Alias            string            `mapstructure:"alias"`
	Description      string            `mapstructure:"description"`
	About            string            `mapstructure:"about"`
	Usage            string            `mapstructure:"usage"`
	Architectures    []string          `mapstructure:"architectures"`
	OperatingSystems []string          `mapstructure:"operatingSystems"`
	Tags             map[string]string `mapstructure:"tags"`
}

type PrivateHelmChart struct {
	Name       string            `mapstructure:"name"`
	Version    string            `mapstructure:"version"`
	ScanOnPush bool              `mapstructure:"scanOnPush"`
	Mutability string            `mapstructure:"mutability"`
	Chart      string            `mapstructure:"chart"`
	Tags       map[string]string `mapstructure:"tags"`
}

type PublicHelmChart struct {
	Name             string            `mapstructure:"name"`
	Version          string            `mapstructure:"version"`
	ScanOnPush       bool              `mapstructure:"scanOnPush"`
	Mutability       string            `mapstructure:"mutability"`
	Chart            string            `mapstructure:"chart"`
	Alias            string            `mapstructure:"alias"`
	Description      string            `mapstructure:"description"`
	About            string            `mapstructure:"about"`
	Usage            string            `mapstructure:"usage"`
	Architectures    []string          `mapstructure:"architectures"`
	OperatingSystems []string          `mapstructure:"operatingSystems"`
	Tags             map[string]string `mapstructure:"tags"`
}

type ContainerRegistryClient interface {
	GetRepositoryPolicy(account string, repository string) *ecr.GetRepositoryPolicyOutput
	AddAccountToPolicy(account, repository, policy string)
	CreateRepositoryPolicy(account, repository, policy string)
	RemoveAccountFromPolicy(account, repository, remove string)
	CreateUpdateHelm(account, region string, hc PrivateHelmChart)
	CreateUpdateDocker(account, region string, r PrivateDockerImage)
	MaybeCreateRepository(account string, repos []PrivateRepository)
	Destroy(account, repository string)
	ConvertDockerImagesToRepositories(images []PrivateDockerImage) []PrivateRepository
	ConvertHelmChartsToRepositories(charts []PrivateHelmChart) []PrivateRepository
}

type PublicClient struct {
	Api *ecrpublic.Client
}

type PrivateClient struct {
	Api *ecr.Client
}

type Public struct {
	Region string              `mapstructure:"region"`
	Images []PublicDockerImage `mapstructure:"images"`
	Charts []PublicHelmChart   `mapstructure:"charts"`
}

type Private struct {
	Region string               `mapstructure:"region"`
	Images []PrivateDockerImage `mapstructure:"images"`
	Charts []PrivateHelmChart   `mapstructure:"charts"`
}

type Prep struct {
	Account      string  `mapstructure:"account"`
	Environment  string  `mapstructure:"environment"`
	Version      string  `mapstructure:"version"`
	Organization string  `mapstructure:"organization"`
	Name         string  `mapstructure:"name"`
	Alias        string  `mapstructure:"alias"`
	Domain       string  `mapstructure:"domain"`
	Private      Private `mapstructure:"private"`
	Public       Public  `mapstructure:"public"`
}

const repositoryPolicy = `
{
	"Version": "2012-10-17",
	"Statement": [ %s ]
}
`

func NewPrivateEcrClient(region string) *PrivateClient {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default private ecr configurations")
		return nil
	}

	api := ecr.NewFromConfig(cfg)

	return &PrivateClient{Api: api}
}

func NewPublicEcrClient(region string) *PublicClient {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error loading default public ecr configurations")
		return nil
	}

	api := ecrpublic.NewFromConfig(cfg)

	return &PublicClient{Api: api}
}

func (client *PrivateClient) CreateUpdateHelmChart(account, region string, hc PrivateHelmChart) bool {
	b := fmt.Sprintf("aws ecr get-login-password --region %s | helm registry login --username AWS --password-stdin %s.dkr.ecr.%s.amazonaws.com", region, account, region)
	c := strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("authenticate helm")

	err := terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error authenticating ecr")
		os.Exit(1)
		return false
	}

	b = fmt.Sprintf("helm package %s", hc.Chart)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("package helm chart")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error building helm chart")
		os.Exit(1)
		return false
	}

	i := strings.LastIndex(hc.Name, "/")
	tar := hc.Name[i+1:]
	repository := hc.Name[:i]

	b = fmt.Sprintf("helm push %s-%s.tgz oci://%s.dkr.ecr.%s.amazonaws.com/%s", tar, hc.Version, account, region, repository)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("push helm chart")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error pushing helm chart")
		os.Exit(1)
		return false
	}

	b = fmt.Sprintf("rm %s-%s.tgz", tar, hc.Version)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("remove helm package")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error deleting generated helm package")
	}

	return true
}

func (client *PublicClient) CreateUpdateHelmChart(alias, region string, hc PublicHelmChart) bool {
	b := fmt.Sprintf("aws ecr-public get-login-password --region %s | helm registry login --username AWS --password-stdin public.ecr.aws/%s", region, alias)
	c := strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("authenticate helm")

	err := terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error authenticating ecr")
		os.Exit(1)
		return false
	}

	b = fmt.Sprintf("helm package %s", hc.Chart)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("package helm chart")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error building helm chart")
		os.Exit(1)
		return false
	}

	i := strings.LastIndex(hc.Name, "/")
	tar := hc.Name[i+1:]
	repository := hc.Name[:i]

	b = fmt.Sprintf("helm push %s-%s.tgz oci://public.ecr.aws/%s/%s", tar, hc.Version, alias, repository)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("push helm chart")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error pushing helm chart")
		os.Exit(1)
		return false
	}

	b = fmt.Sprintf("rm %s-%s.tgz", tar, hc.Version)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("remove helm package")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error deleting generated helm package")
	}

	return true
}

func (client *PrivateClient) CreateUpdateDockerImage(account, region string, r PrivateDockerImage) bool {
	b := fmt.Sprintf("aws ecr get-login-password --region %s | docker login --username AWS --password-stdin %s.dkr.ecr.%s.amazonaws.com", region, account, region)
	c := strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("authenticate ecr")

	err := terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error authenticating ecr")
		os.Exit(1)
		return false
	}

	tags := fmt.Sprintf("-t %s.dkr.ecr.%s.amazonaws.com/%s:%s", account, region, r.Name, r.Version)
	tags = fmt.Sprintf("%s -t %s.dkr.ecr.%s.amazonaws.com/%s:%s", tags, account, region, r.Name, time.Now().Format("20060102"))
	tags = fmt.Sprintf("%s -t %s.dkr.ecr.%s.amazonaws.com/%s:latest", tags, account, region, r.Name)
	b = fmt.Sprintf("docker buildx build --provenance=false --platform linux/amd64 -f %s %s %s --push", r.Dockerfile, tags, r.Context)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("build docker container")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error building docker image")
		os.Exit(1)
		return false
	}

	return true
}

func (client *PublicClient) CreateUpdateDockerImage(alias, region string, r PublicDockerImage) bool {
	b := fmt.Sprintf("aws ecr-public get-login-password --region %s | docker login --username AWS --password-stdin public.ecr.aws/%s", region, alias)
	c := strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("authenticate ecr")

	err := terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error authenticating ecr")
		os.Exit(1)
		return false
	}

	tags := fmt.Sprintf("-t public.ecr.aws/%s/%s:%s", r.Alias, r.Name, r.Version)
	tags = fmt.Sprintf("%s -t public.ecr.aws/%s/%s:%s", tags, r.Alias, r.Name, time.Now().Format("20060102"))
	tags = fmt.Sprintf("%s -t public.ecr.aws/%s/%s:latest", tags, r.Alias, r.Name)
	b = fmt.Sprintf("docker buildx build --provenance=false --platform linux/amd64 -f %s %s %s", r.Dockerfile, tags, r.Context)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("build docker container")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error building docker image")
		os.Exit(1)
		return false
	}

	b = fmt.Sprintf("docker push -a public.ecr.aws/%s/%s", r.Alias, r.Name)
	c = strings.Join(strings.Fields(b), " ")

	logger.Logger.Info().
		Str("command", b).
		Msg("push docker container")

	err = terminal.ExecuteCommand(c, time.Duration(5)*time.Minute)
	if err != nil {
		logger.Logger.Err(err).Msg("error pushing docker image")
		os.Exit(1)
		return false
	}

	return true
}
