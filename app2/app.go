package app2

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kudrykv/latex-yearly-planner/app2/planners"
	"github.com/kudrykv/latex-yearly-planner/app2/planners/breadcrumb"
	"github.com/kudrykv/latex-yearly-planner/app2/planners/common"
	"github.com/kudrykv/latex-yearly-planner/app2/planners/mos"
	"github.com/kudrykv/latex-yearly-planner/app2/planners2"
	mos2 "github.com/kudrykv/latex-yearly-planner/app2/planners2/mos"
	"github.com/kudrykv/latex-yearly-planner/app2/types"
	"github.com/urfave/cli/v2"
)

type App struct {
	app *cli.App
}

const (
	yearFlag     = "year"
	sectionsFlag = "sections"
	weekdayFlag  = "weekday"
	framesFlag   = "show-frames"
	linksFlag    = "show-links"

	deviceNameFlag = "device-name"

	handFlag = "hand"

	layoutPathFlag     = "layout-path"
	parametersPathFlag = "parameters-path"
)

func New(reader io.Reader, writer, errWriter io.Writer) *App {
	return (&App{}).
		setupCli(reader, writer, errWriter)
}

func (r *App) setupCli(reader io.Reader, writer, errWriter io.Writer) *App {
	r.app = &cli.App{
		Name: "plannergen",

		Reader:    reader,
		Writer:    writer,
		ErrWriter: errWriter,

		Commands: cli.Commands{
			&cli.Command{
				Name: "template",
				Subcommands: cli.Commands{
					templateCommand[mos.UI](planners.MonthsOnSidesTemplate),
					templateCommand[breadcrumb.UI](planners.BreadcrumbTemplate),

					&cli.Command{
						Name: "mos2",

						Flags: flags(),

						Action: func(appContext *cli.Context) error {
							var layout types.Layout
							if err := readToml(appContext.String(layoutPathFlag), &layout); err != nil {
								return fmt.Errorf("read layout: %w", err)
							}

							var parameters mos2.Parameters
							if err := readToml(appContext.String(parametersPathFlag), &parameters); err != nil {
								return fmt.Errorf("read parameters: %w", err)
							}

							layout.Misc = parameters

							planner, err := planners2.New("mos", layout)
							if err != nil {
								return fmt.Errorf("new planner: %w", err)
							}

							if err := planner.Generate(); err != nil {
								return fmt.Errorf("generate: %w", err)
							}

							if err := planner.WriteTeXTo("./out"); err != nil {
								return fmt.Errorf("write tex: %w", err)
							}

							if err := planner.Compile(appContext.Context); err != nil {
								return fmt.Errorf("compile: %w", err)
							}

							return nil
						},
					},
				},
			},
		},
	}

	return r
}

func readToml(path string, dst any) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if err := toml.Unmarshal(fileBytes, dst); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}

func templateCommand[T any](template string) *cli.Command {
	return &cli.Command{
		Name: template,

		Flags: flags(),

		Action: func(appContext *cli.Context) error {
			hand := common.RightHand
			if appContext.String(handFlag) == "left" {
				hand = common.LeftHand
			}

			var ui T
			if path := appContext.String(layoutPathFlag); path != "" {
				fileBytes, err := os.ReadFile(appContext.String(layoutPathFlag))
				if err != nil {
					return fmt.Errorf("read file: %w", err)
				}

				if err = toml.Unmarshal(fileBytes, &ui); err != nil {
					return fmt.Errorf("unmarshal ui: %w", err)
				}
			}

			params := common.NewParams(
				common.ParamWithYear[T](appContext.Int(yearFlag)),
				common.ParamWithDeviceName[T](appContext.String(deviceNameFlag)),
				common.ParamWithSections[T](appContext.StringSlice(sectionsFlag)),
				common.ParamWithWeekday[T](time.Weekday(appContext.Int(weekdayFlag))),
				common.ParamWithMainHand[T](hand),
				common.ParamWithFrames[T](appContext.Bool(framesFlag)),
				common.ParamWithLinks[T](appContext.Bool(linksFlag)),
				common.ParamWithUI[T](ui),
			)

			planner, err := planners.New(template, params)
			if err != nil {
				return fmt.Errorf("new planner: %w", err)
			}

			if err = planner.Generate(); err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			if err = planner.WriteTeXTo("./out"); err != nil {
				return fmt.Errorf("write to ./out: %w", err)
			}

			if err = planner.Compile(appContext.Context); err != nil {
				return fmt.Errorf("compile: %w", err)
			}

			return nil
		},
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: layoutPathFlag},
		&cli.StringFlag{Name: parametersPathFlag},
	}
}

func (r *App) Run(args []string) error {
	return r.app.Run(args)
}
