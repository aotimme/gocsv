package main

import (
  "encoding/csv"
  "errors"
  "flag"
  "io"
  "os"
)

func StackFiles(inreaders []io.Reader, groupName string, groups []string)  {
  shouldAppendGroup := groupName != ""
  writer := csv.NewWriter(os.Stdout)
  readers := make([]*csv.Reader, len(inreaders))
  for i, inreader := range inreaders {
    readers[i] = csv.NewReader(inreader)
  }

  // Check that the headers match
  headers := make([][]string, len(inreaders))
  for i, reader := range readers {
    header, err := reader.Read()
    if err != nil {
      panic(err)
    }
    headers[i] = header
  }
  firstHeader := headers[0]
  for i, header := range headers {
    if i == 0 {
      continue
    }
    if len(firstHeader) != len(header) {
      panic(errors.New("Headers do not match"))
    }
    for j, elem := range firstHeader {
      if elem != header[j] {
        panic(errors.New("Headers do not match"))
      }
    }
  }
  if shouldAppendGroup {
    firstHeader = append(firstHeader, groupName)
  }
  writer.Write(firstHeader)
  writer.Flush()

  // Go through the files
  for i, reader := range readers {
    for {
      row, err := reader.Read()
      if err != nil {
        if err == io.EOF {
          break
        } else {
          panic(err)
        }
      }
      if shouldAppendGroup {
        row = append(row, groups[i])
      }
      writer.Write(row)
      writer.Flush()
    }
  }
}


func RunStack(args []string) {
  fs := flag.NewFlagSet("stack", flag.ExitOnError)
  var groupName, groupsString string
  var useFilenames bool
  fs.StringVar(&groupName, "group-name", "", "Name of the column for grouping")
  fs.StringVar(&groupsString, "groups", "", "Group to display for each file")
  fs.BoolVar(&useFilenames, "filenames", false, "Use the filename for groups")
  err := fs.Parse(args)
  if err != nil {
    panic(err)
  }
  filenames := fs.Args()

  hasSpecifiedGroups := groupsString != ""
  if hasSpecifiedGroups && useFilenames {
    panic(errors.New("Cannot specify both --filename and --groups"))
  }

  shouldAppendGroup := hasSpecifiedGroups || useFilenames

  var groups []string
  if hasSpecifiedGroups {
    groups = GetArrayFromCsvString(groupsString)
  } else if useFilenames {
    groups = filenames
  }

  if shouldAppendGroup && len(filenames) != len(groups) {
    panic(errors.New("Number of files and groups are not equal"))
  }

  var groupColumnName string
  if (groupName != "") {
    groupColumnName = groupName
  } else if useFilenames {
    groupColumnName = "File"
  } else if shouldAppendGroup {
    groupColumnName = "Group"
  } else {
    groupColumnName = ""
  }

  inreaders := make([]io.Reader, len(filenames))
  for i, filename := range filenames {
    file, err := os.Open(filename)
    if err != nil {
      panic(err)
    }
    defer file.Close()
    inreaders[i] = file
  }
  StackFiles(inreaders, groupColumnName, groups)
}
