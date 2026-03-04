package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/usama/linkmngr-cli/internal/client"
	"github.com/usama/linkmngr-cli/internal/config"
)

var Version = "dev"
var outputFormat = "table"

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "linkmngr",
		Short:         "CLI for LinkMngr",
		Example:       "  linkmngr auth login <token>\n  linkmngr link list --page 1\n  linkmngr api request GET /links",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (`table` or `json`)")

	root.AddCommand(newAuthCmd())
	root.AddCommand(newLinksCmd())
	root.AddCommand(newBrandsCmd())
	root.AddCommand(newAnalyticsCmd())
	root.AddCommand(newDomainsCmd())
	root.AddCommand(newPagesCmd())
	root.AddCommand(newAPICmd())
	root.AddCommand(newVersionCmd())

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Example: "  linkmngr version",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), Version)
		},
	}
}

func newAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Example: "  linkmngr auth login <token>\n  linkmngr auth status\n  linkmngr auth logout",
	}

	authCmd.AddCommand(
		&cobra.Command{
			Use:     "login <token>",
			Aliases: []string{"set-token"},
			Short:   "Authenticate with an API token",
			Example: "  linkmngr auth login <token>",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return err
				}
				cfg.Token = args[0]
				if err := config.Save(cfg); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Token saved.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "set-base-url <url>",
			Short: "Set the API base URL in local config",
			Example: "  linkmngr auth set-base-url https://api.linkmngr.com/v1",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := config.Load()
				if err != nil {
					return err
				}
				cfg.BaseURL = strings.TrimRight(args[0], "/")
				if err := config.Save(cfg); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Base URL saved.")
				return nil
			},
		},
		&cobra.Command{
			Use:     "status",
			Aliases: []string{"whoami"},
			Short:   "Show authentication status",
			Example: "  linkmngr auth status",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := apiClientFromConfig()
				if err != nil {
					return err
				}
				user, err := c.WhoAmI(cmd.Context())
				if err != nil {
					return err
				}
				return printOutput(cmd, user)
			},
		},
		&cobra.Command{
			Use:     "logout",
			Aliases: []string{"revoke"},
			Short:   "Revoke current token",
			Example: "  linkmngr auth logout",
			RunE: func(cmd *cobra.Command, _ []string) error {
				c, err := apiClientFromConfig()
				if err != nil {
					return err
				}
				if err := c.RevokeToken(cmd.Context()); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Token revoked.")
				return nil
			},
		},
	)

	return authCmd
}

func newLinksCmd() *cobra.Command {
	linksCmd := &cobra.Command{
		Use:     "link",
		Aliases: []string{"links"},
		Short:   "Manage links",
		Example: "  linkmngr link list\n  linkmngr link get 123\n  linkmngr link create https://example.com",
	}

	linksCmd.AddCommand(newLinksListCmd())
	linksCmd.AddCommand(newLinksGetCmd())
	linksCmd.AddCommand(newLinksCreateCmd())
	linksCmd.AddCommand(newLinksStatsCmd())

	return linksCmd
}

func newLinksListCmd() *cobra.Command {
	var page int
	var brandID int
	var domain string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List links",
		Example: "  linkmngr link list\n  linkmngr link list --page 2 --brand-id 12",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			results, err := c.ListLinks(cmd.Context(), client.ListLinksParams{
				Page:    page,
				BrandID: brandID,
				Domain:  domain,
			})
			if err != nil {
				return err
			}
			return printOutput(cmd, results)
		},
	}

	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number (1-based)")
	cmd.Flags().IntVar(&brandID, "brand-id", 0, "Filter by brand ID")
	cmd.Flags().StringVar(&domain, "domain", "", "Filter by domain")
	return cmd
}

func newLinksGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <link-id>",
		Aliases: []string{"view"},
		Short:   "Get a link",
		Example: "  linkmngr link get 123\n  linkmngr link view 123",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0], "link ID")
			if err != nil {
				return err
			}
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			link, err := c.GetLink(cmd.Context(), id)
			if err != nil {
				return err
			}
			return printOutput(cmd, link)
		},
	}
}

func newLinksCreateCmd() *cobra.Command {
	var domain string
	var slug string
	var brandID int

	cmd := &cobra.Command{
		Use:   "create <destination>",
		Short: "Create a link",
		Example: "  linkmngr link create https://example.com\n  linkmngr link create https://example.com --domain linkmn.gr --slug launch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			in := client.CreateLinkRequest{
				Destination: args[0],
				Domain:      domain,
				Slug:        slug,
				BrandID:     brandID,
			}
			link, err := c.CreateLink(cmd.Context(), in)
			if err != nil {
				return err
			}
			return printOutput(cmd, link)
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "Short domain to use")
	cmd.Flags().StringVar(&slug, "slug", "", "Preferred slug")
	cmd.Flags().IntVar(&brandID, "brand-id", 0, "Attach to brand ID")
	return cmd
}

func newLinksStatsCmd() *cobra.Command {
	var start string
	var end string
	var timeUnit string
	var groupBy string

	cmd := &cobra.Command{
		Use:   "stats <link-id>",
		Short: "Get link stats",
		Example: "  linkmngr link stats 123 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0], "link ID")
			if err != nil {
				return err
			}
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			stats, err := c.GetLinkStats(cmd.Context(), id, client.StatsParams{
				Start:    start,
				End:      end,
				TimeUnit: timeUnit,
				GroupBy:  groupBy,
			})
			if err != nil {
				return err
			}
			return printOutput(cmd, stats)
		},
	}

	cmd.Flags().StringVar(&start, "start", "", "Start time in ISO 8601 format")
	cmd.Flags().StringVar(&end, "end", "", "End time in ISO 8601 format")
	cmd.Flags().StringVar(&timeUnit, "time-unit", "day", "Time unit: hour|day|week|month|year")
	cmd.Flags().StringVar(&groupBy, "group-by", "", "Group by: device|device_type|country|browser|platform|referrer")
	_ = cmd.MarkFlagRequired("start")
	_ = cmd.MarkFlagRequired("end")
	return cmd
}

func newBrandsCmd() *cobra.Command {
	brandsCmd := &cobra.Command{
		Use:     "brand",
		Aliases: []string{"brands"},
		Short:   "Manage brands",
		Example: "  linkmngr brand list\n  linkmngr brand get 12\n  linkmngr brand domain-check 12 linkmn.gr",
	}

	brandsCmd.AddCommand(newBrandsListCmd())
	brandsCmd.AddCommand(newBrandsGetCmd())
	brandsCmd.AddCommand(newBrandsCheckDomainCmd())

	return brandsCmd
}

func newBrandsListCmd() *cobra.Command {
	var page int
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List brands",
		Example: "  linkmngr brand list\n  linkmngr brand list --page 2",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			results, err := c.ListBrands(cmd.Context(), client.ListBrandsParams{Page: page})
			if err != nil {
				return err
			}
			return printOutput(cmd, results)
		},
	}
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number (1-based)")
	return cmd
}

func newBrandsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <brand-id>",
		Aliases: []string{"view"},
		Short:   "Get a brand",
		Example: "  linkmngr brand get 12\n  linkmngr brand view 12",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0], "brand ID")
			if err != nil {
				return err
			}
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			brand, err := c.GetBrand(cmd.Context(), id)
			if err != nil {
				return err
			}
			return printOutput(cmd, brand)
		},
	}
}

func newBrandsCheckDomainCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "domain-check <brand-id> <domain>",
		Aliases: []string{"check-domain"},
		Short:   "Check brand domain configuration",
		Example: "  linkmngr brand domain-check 12 linkmn.gr",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0], "brand ID")
			if err != nil {
				return err
			}
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			result, err := c.CheckBrandDomain(cmd.Context(), id, args[1])
			if err != nil {
				return err
			}
			return printOutput(cmd, result)
		},
	}
}

func newAnalyticsCmd() *cobra.Command {
	var brandID int
	var start string
	var end string
	var timeUnit string
	var groupBy string

	cmd := &cobra.Command{
		Use:   "analytics",
		Short: "Get metrics about your links",
		Example: "  linkmngr analytics --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00 --group-by country",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			result, err := c.GetAnalytics(cmd.Context(), client.StatsParams{
				BrandID:  brandID,
				Start:    start,
				End:      end,
				TimeUnit: timeUnit,
				GroupBy:  groupBy,
			})
			if err != nil {
				return err
			}
			return printOutput(cmd, result)
		},
	}

	cmd.Flags().IntVar(&brandID, "brand-id", 0, "Filter by brand ID")
	cmd.Flags().StringVar(&start, "start", "", "Start time in ISO 8601 format")
	cmd.Flags().StringVar(&end, "end", "", "End time in ISO 8601 format")
	cmd.Flags().StringVar(&timeUnit, "time-unit", "day", "Time unit: hour|day|week|month|year")
	cmd.Flags().StringVar(&groupBy, "group-by", "", "Group by: device|device_type|country|browser|platform|referrer")
	_ = cmd.MarkFlagRequired("start")
	_ = cmd.MarkFlagRequired("end")
	return cmd
}

func newDomainsCmd() *cobra.Command {
	domainsCmd := &cobra.Command{
		Use:     "domain",
		Aliases: []string{"domains"},
		Short:   "Manage domains",
		Example: "  linkmngr domain list",
	}
	domainsCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available domains",
		Example: "  linkmngr domain list",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			domains, err := c.ListDomains(cmd.Context())
			if err != nil {
				return err
			}
			return printOutput(cmd, domains)
		},
	})
	return domainsCmd
}

func newAPICmd() *cobra.Command {
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Make raw API requests for advanced/undocumented endpoints",
		Example: "  linkmngr api request GET /links\n  linkmngr api request POST /links --set destination=https://example.com",
	}
	apiCmd.AddCommand(newAPIRequestCmd())
	return apiCmd
}

func newPagesCmd() *cobra.Command {
	pagesCmd := &cobra.Command{
		Use:     "page",
		Aliases: []string{"pages"},
		Short:   "Manage bio pages",
		Example: "  linkmngr page list\n  linkmngr page get 44\n  linkmngr page stats 44 --start <iso> --end <iso>",
	}
	pagesCmd.AddCommand(newPagesListCmd())
	pagesCmd.AddCommand(newPagesGetCmd())
	pagesCmd.AddCommand(newPagesStatsCmd())
	pagesCmd.AddCommand(newPagesHitsCmd())
	return pagesCmd
}

func newPagesListCmd() *cobra.Command {
	var page int
	var brandID int
	var domain string
	var customDomainID int
	var slug string
	var query string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List pages",
		Example: "  linkmngr page list --brand-id 12 --search launch",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			results, err := c.ListPages(cmd.Context(), client.ListPagesParams{
				Page:           page,
				BrandID:        brandID,
				Domain:         domain,
				CustomDomainID: customDomainID,
				Slug:           slug,
				Query:          query,
			})
			if err != nil {
				return err
			}
			return printOutput(cmd, results)
		},
	}

	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number (1-based)")
	cmd.Flags().IntVar(&brandID, "brand-id", 0, "Filter by brand ID")
	cmd.Flags().StringVar(&domain, "domain", "", "Filter by domain")
	cmd.Flags().IntVar(&customDomainID, "custom-domain-id", 0, "Filter by custom domain ID")
	cmd.Flags().StringVar(&slug, "slug", "", "Filter by slug")
	cmd.Flags().StringVar(&query, "search", "", "Search by title/description")
	return cmd
}

func newPagesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <page-id>",
		Aliases: []string{"view"},
		Short:   "Get a page",
		Example: "  linkmngr page get 44\n  linkmngr page view 44",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0], "page ID")
			if err != nil {
				return err
			}
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			page, err := c.GetPage(cmd.Context(), id)
			if err != nil {
				return err
			}
			return printOutput(cmd, page)
		},
	}
}

func newPagesStatsCmd() *cobra.Command {
	var start string
	var end string
	var timeUnit string
	var groupBy string

	cmd := &cobra.Command{
		Use:   "stats <page-id>",
		Short: "Get page analytics stats",
		Example: "  linkmngr page stats 44 --start 2026-03-01T00:00:00+00:00 --end 2026-03-03T00:00:00+00:00",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0], "page ID")
			if err != nil {
				return err
			}
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			stats, err := c.GetPageStats(cmd.Context(), id, client.StatsParams{
				Start:    start,
				End:      end,
				TimeUnit: timeUnit,
				GroupBy:  groupBy,
			})
			if err != nil {
				return err
			}
			return printOutput(cmd, stats)
		},
	}

	cmd.Flags().StringVar(&start, "start", "", "Start time in ISO 8601 format")
	cmd.Flags().StringVar(&end, "end", "", "End time in ISO 8601 format")
	cmd.Flags().StringVar(&timeUnit, "time-unit", "day", "Time unit: hour|day|week|month|year")
	cmd.Flags().StringVar(&groupBy, "group-by", "", "Group by: device|device_type|country|browser|platform|referrer")
	_ = cmd.MarkFlagRequired("start")
	_ = cmd.MarkFlagRequired("end")
	return cmd
}

func newPagesHitsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hits <page-id>",
		Short: "Get recent page hits",
		Example: "  linkmngr page hits 44",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0], "page ID")
			if err != nil {
				return err
			}
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}
			hits, err := c.GetPageHits(cmd.Context(), id)
			if err != nil {
				return err
			}
			return printOutput(cmd, hits)
		},
	}
}

func newAPIRequestCmd() *cobra.Command {
	var jsonBody string
	var jsonBodyFile string
	var setPairs []string

	cmd := &cobra.Command{
		Use:   "request <method> <path>",
		Short: "Send a raw request to LinkMngr API (e.g. /links, /brands, /pages)",
		Example: "  linkmngr api request GET /links\n  linkmngr api request PATCH /brands/12 --data-file payload.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := apiClientFromConfig()
			if err != nil {
				return err
			}

			body, err := buildRequestBody(jsonBody, jsonBodyFile, setPairs)
			if err != nil {
				return err
			}

			resp, err := c.RawRequest(cmd.Context(), args[0], args[1], body)
			if err != nil {
				return err
			}
			if resp == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "Request completed.")
				return nil
			}
			return printOutput(cmd, resp)
		},
	}

	cmd.Flags().StringVar(&jsonBody, "data", "", "Raw JSON request body")
	cmd.Flags().StringVar(&jsonBodyFile, "data-file", "", "Path to JSON file body")
	cmd.Flags().StringSliceVar(&setPairs, "set", nil, "Set body field as key=value (repeatable)")
	return cmd
}

func apiClientFromConfig() (*client.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return apiClient(cfg)
}

func apiClient(cfg config.Config) (*client.Client, error) {
	if cfg.Token == "" {
		return nil, errors.New("missing API token; set LINKMNGR_TOKEN or run `linkmngr auth login <token>`")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.linkmngr.com/v1"
	}
	return client.New(cfg.BaseURL, cfg.Token), nil
}

func parseID(raw string, field string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: must be an integer (try a numeric ID like `123`)", field, raw)
	}
	if id <= 0 {
		return 0, fmt.Errorf("invalid %s %q: must be > 0 (try a numeric ID like `123`)", field, raw)
	}
	return id, nil
}

func buildRequestBody(jsonBody string, jsonBodyFile string, setPairs []string) (map[string]any, error) {
	if jsonBody != "" && jsonBodyFile != "" {
		return nil, errors.New("use either --data or --data-file, not both (try: `linkmngr api request POST /links --set destination=https://example.com`)")
	}

	if jsonBodyFile != "" {
		b, err := os.ReadFile(jsonBodyFile)
		if err != nil {
			return nil, fmt.Errorf("read data file: %w", err)
		}
		jsonBody = string(b)
	}

	if strings.TrimSpace(jsonBody) != "" {
		var body map[string]any
		if err := json.Unmarshal([]byte(jsonBody), &body); err != nil {
			return nil, fmt.Errorf("invalid JSON body: %w", err)
		}
		return body, nil
	}

	if len(setPairs) == 0 {
		return nil, nil
	}
	body := make(map[string]any, len(setPairs))
	for _, p := range setPairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("invalid --set value %q: expected key=value (try: --set destination=https://example.com)", p)
		}
		k = strings.TrimSpace(k)
		if k == "" {
			return nil, fmt.Errorf("invalid --set value %q: empty key", p)
		}
		body[k] = strings.TrimSpace(v)
	}
	return body, nil
}

func printOutput(cmd *cobra.Command, v any) error {
	switch strings.ToLower(strings.TrimSpace(outputFormat)) {
	case "", "json":
		return printJSON(cmd, v)
	case "table":
		return printTable(cmd, v)
	default:
		return fmt.Errorf("invalid output format %q: use `--output table` or `--output json`", outputFormat)
	}
}

func printTable(cmd *cobra.Command, v any) error {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	defer w.Flush()

	switch t := v.(type) {
	case client.PaginatedResults[client.Link]:
		fmt.Fprintln(w, "ID\tSHORT_LINK\tDESTINATION\tDOMAIN\tCLICKS\tUNIQUE_CLICKS\tCREATED_AT")
		for _, it := range t.Items {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%d\t%s\n", it.ID, it.Link, it.Destination, it.Domain, it.Clicks, it.UniqueClicks, it.CreatedAt)
		}
		return nil
	case client.PaginatedResults[client.Brand]:
		fmt.Fprintln(w, "ID\tNAME\tDOMAINS")
		for _, it := range t.Items {
			domains := make([]string, 0, len(it.Domains))
			for _, d := range it.Domains {
				if d.Domain != "" {
					domains = append(domains, d.Domain)
				}
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", it.ID, it.Name, strings.Join(domains, ","))
		}
		return nil
	case client.PaginatedResults[client.Page]:
		fmt.Fprintln(w, "ID\tTITLE\tSLUG\tURL\tCLICKS\tUNIQUE_CLICKS\tLAST_CLICK_AT")
		for _, it := range t.Items {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%d\t%s\n", it.ID, it.Title, it.Slug, it.URL, it.Clicks, it.UniqueClicks, it.LastClickAt)
		}
		return nil
	case client.Link:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "id\t%d\n", t.ID)
		fmt.Fprintf(w, "short_link\t%s\n", t.Link)
		fmt.Fprintf(w, "destination\t%s\n", t.Destination)
		fmt.Fprintf(w, "domain\t%s\n", t.Domain)
		fmt.Fprintf(w, "slug\t%s\n", t.Slug)
		fmt.Fprintf(w, "brand_id\t%d\n", t.BrandID)
		fmt.Fprintf(w, "clicks\t%d\n", t.Clicks)
		fmt.Fprintf(w, "unique_clicks\t%d\n", t.UniqueClicks)
		fmt.Fprintf(w, "created_at\t%s\n", t.CreatedAt)
		return nil
	case client.Brand:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "id\t%d\n", t.ID)
		fmt.Fprintf(w, "name\t%s\n", t.Name)
		domains := make([]string, 0, len(t.Domains))
		for _, d := range t.Domains {
			if d.Domain != "" {
				domains = append(domains, d.Domain)
			}
		}
		fmt.Fprintf(w, "domains\t%s\n", strings.Join(domains, ","))
		return nil
	case client.Page:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "id\t%d\n", t.ID)
		fmt.Fprintf(w, "title\t%s\n", t.Title)
		fmt.Fprintf(w, "slug\t%s\n", t.Slug)
		fmt.Fprintf(w, "url\t%s\n", t.URL)
		fmt.Fprintf(w, "description\t%s\n", t.Description)
		fmt.Fprintf(w, "brand_id\t%d\n", t.BrandID)
		fmt.Fprintf(w, "custom_domain_id\t%d\n", t.CustomDomainID)
		fmt.Fprintf(w, "clicks\t%d\n", t.Clicks)
		fmt.Fprintf(w, "unique_clicks\t%d\n", t.UniqueClicks)
		fmt.Fprintf(w, "last_click_at\t%s\n", t.LastClickAt)
		return nil
	case client.BrandDomain:
		fmt.Fprintln(w, "DOMAIN\tACTIVE\tTYPE")
		fmt.Fprintf(w, "%s\t%t\t%s\n", t.Domain, t.Active, t.Type)
		return nil
	case []string:
		fmt.Fprintln(w, "DOMAIN")
		for _, d := range t {
			fmt.Fprintln(w, d)
		}
		return nil
	case map[string]any:
		return printMapTable(w, t)
	case []map[string]any:
		for i, itm := range t {
			if i == 0 {
				fmt.Fprintln(w, "INDEX\tDATA")
			}
			b, err := json.Marshal(itm)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "%d\t%s\n", i+1, string(b))
		}
		return nil
	default:
		return printJSON(cmd, v)
	}
}

func printMapTable(w *tabwriter.Writer, m map[string]any) error {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintln(w, "FIELD\tVALUE")
	for _, k := range keys {
		b, err := json.Marshal(m[k])
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s\t%s\n", k, string(b))
	}
	return nil
}

func printJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

func Run(ctx context.Context, args []string) error {
	root := NewRootCmd()
	root.SetArgs(args)
	root.SetContext(ctx)
	return root.Execute()
}
