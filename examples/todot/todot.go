package main

import (
	"../../src/graph/_obj/graph"
	
	"fmt"
	"flag"
	"os"
	"path"

	"github.com/StepLg/go-erx/src/erx"
)

func main() {
    defer func() {
        if err := recover(); err != nil {
            if errErx, ok := err.(erx.Error); ok {
                formatter := erx.NewStringFormatter("  ")
                fmt.Println(formatter.Format(errErx))
            }
        }
    }()

	flag_help := flag.Bool("help", false, "Display this help.")
	flag_inputFile := flag.String("in", "", 
`Input file with messages. If file extension is .ugr, .dgr or .mgr then
graph type automaticly set to undirected, directed or mixed respectively.
If flag doesn't set, then read from stdin.`)
	flag_outputFile := flag.String("out", "",
`Output file. If it isn't set and flag -autoname isn't set too, then
output to stdout.`)
	flag_autoname := flag.Bool("autoname", false,
`Generate output file name automaticly from input file name with 
replacing it's extension to ".dot."`)
	flag_type := flag.String("type", "", 
`[u|d|m] -- Graph type: undirected, directed or mixed respectively.`)
	
	flag.Parse()

	if *flag_help {
		flag.PrintDefaults()
		return
	}
	
	infile := os.Stdin
	if *flag_inputFile!="" {
		var err os.Error
		infile, err = os.Open(*flag_inputFile, os.O_RDONLY, 0000)
		if err!=nil {
			erxErr := erx.NewSequent("Can't open input file.", err)
			erxErr.AddV("file name", *flag_inputFile)
			panic(erxErr)
		}
		
		infileExt := path.Ext(*flag_inputFile)
		
		if *flag_autoname {
			// generating autoname for output file
			fpath, fname := path.Split(*flag_inputFile)
			baseFileName := fname
			baseFileName = fname[0:len(fname)-len(infileExt)]
			
			*flag_outputFile = path.Join(fpath, baseFileName + ".dot")
		}
		
		// initializing graph type from known file extensions.
		switch infileExt {
			case ".ugr":
				*flag_type = "u"
			case ".dgr":
				*flag_type = "d"
			case ".mgr":
				*flag_type = "m"
		}
	}
	
	outfile := os.Stdout
	if *flag_outputFile!="" {
		var err os.Error
		outfile, err = os.Open(*flag_outputFile, os.O_WRONLY | os.O_CREAT | os.O_TRUNC, 0644)
		if err!=nil {
			erxErr := erx.NewSequent("Can't open output file.", err)
			erxErr.AddV("file name", *flag_outputFile)
			panic(erxErr)
		}
	}
	
	switch *flag_type {
		case "u":
			gr := graph.NewUndirectedMap()
			graph.ReadUgraphFile(infile, gr)
			graph.PlotUgraphToDot(gr, outfile, nil, nil)
		case "d":
			gr := graph.NewDirectedMap()
			graph.ReadDgraphFile(infile, gr)
			graph.PlotDgraphToDot(gr, outfile, nil, nil)
		case "m":
			gr := graph.NewMixedMap()
			graph.ReadMgraphFile(infile, gr)
			graph.PlotMgraphToDot(gr, outfile, nil, nil)
		default:
			err := erx.NewError("Unknown type flag.")
			err.AddV("flag value", *flag_type)
			panic(err)
	}
	
	outfile.Close()
	return
}
