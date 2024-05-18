package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Translation Tool"
	app.Usage = "Translate text from one language to another"
	app.Version = "1.0"

	// Define command-line flags and options
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "engine, e",
            Value: "google",
			Usage: "Engine to use",
		},
		cli.StringFlag{
			Name:  "source, s",
                        Value: "auto",
			Usage: "Source Language to translate",
		},
		cli.StringFlag{
			Name:  "target, t",
			Usage: "Target Language to translate",
		},
		cli.StringFlag{
			Name:  "text, txt",
			Usage: "Text to translate",
		},
		cli.BoolFlag{
			Name:  "list-languages, ll",
			Usage: "List Languages supported",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Path to a file (.txt) to translate",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output file for translation result",
		},
	}

	app.Action = func(c *cli.Context) error {
        urls := []string{"https://translate.projectsegfau.lt", "https://mozhi.aryak.me"}
		engine := c.String("engine")
		source := c.String("source")
		target := c.String("target")
		text := c.String("text")

		translate := NewTranslator(urls, engine)

		if c.Bool("list-languages") {
			lang, err := translate.Languages()
            if err != nil {
               panic(err)    
            }
			fmt.Printf("%-25s %-25s\n", "Name", "Code")
			for key, value := range lang {
				fmt.Printf("%-25s %-25s\n", key, value)
			}
		} else if file := c.String("file"); file != "" {
			if output := c.String("output"); output != "" {
				translateAndSaveToFile(translate, source, target, file, output)
			} else {
				result, err := translateFromFile(translate, source, target, file)
                if err != nil {
                   panic(err)    
                }
				fmt.Println(result)
			}
		} else if text != ""{
            if output := c.String("output"); output != "" {
			translateAndSaveToFile(translate, source, target, text, output)
			fmt.Println("Translation result saved in", output)
            } else {
			result, err := translate.Translate(source, target, text)
            if err != nil {
               panic(err)    
            }
			fmt.Println(result)
		}
		} 

		return nil
	}
    if len(os.Args) > 1 {
        //fmt.Println(os.Args)
	    err := app.Run(os.Args)
	    if err != nil {
		    fmt.Println("Error:", err)
	    }
    }else {
        arg := append(os.Args, "help")
        err := app.Run(arg)
	    if err != nil {
		    fmt.Println("Error:", err)
	    }
    }
}

func translateFromFile(translate *Translator, source, target, file string) (string, error) {
    if file != "" {
        text, err := readFile(file)
        if err != nil {
            return "", err // Return the error
        }
        result, err := translate.Translate(source, target, text)
        if err != nil {
           panic(err)    
        }
        return result, nil
    }
    
    return "", fmt.Errorf("File not provided") // Return an error indicating that the file is not provided
}


func translateAndSaveToFile(translate *Translator, source, target, file, output string) (string, error) {
	if file != "" {
		text, err := readFile(file)
		if err != nil {
			return "", fmt.Errorf("Error reading file:", err)
		}
		result, err := translate.Translate(source, target, text)
        if err != nil {
            panic(err)    
        }
		err = writeFile(output, result)
		if err != nil {
			return "", fmt.Errorf("Error writing to output file:", err)
		}
	} 
        return "", fmt.Errorf("File not provided")
    
}



func readFile(filename string) (string, error) {
    file, err := os.ReadFile(filename)
    if err != nil {
        panic(err)    
    }
	return string(file), nil
}

func writeFile(filename, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
    if err != nil {
        panic(err)    
    }
	return nil
}



