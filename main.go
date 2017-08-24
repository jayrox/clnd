package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	flagDebug  = flag.Bool("d", true, "show debug output")
	flagLog    = flag.String("log", "cwd", "directory to write the log. default is current working directory (cwd)")
	flagNoLogo = flag.Bool("nologo", false, "hide the logo, useful for automation logging.")
	f          fFlags
)

func main() {
	flag.Parse()
	// Print the logo :P
	printLogo()

	// Root folder to write logs
	fpLAbs, _ := filepath.Abs(flagString(flagLog))
	f.Log = fpLAbs
	if flagString(flagLog) == "cwd" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		check(err)
		f.Log = dir
	}

	fmt.Println("_____________________")

	// Radarr
	if os.Getenv("radarr_eventtype") == "Download" {
		radarrOnDownload()
	}

	// Sonarr
	if os.Getenv("sonarr_eventtype") == "Download" {
		sonarrOnDownload()
	}

	// Log
	fileName := "[SR] log.log"
	filePath := filepath.Join(f.Log, fileName)

	t := time.Now()
	tn := t.Format("2006-01-02:15:04:05")

	writeToFile(filePath, []byte(tn+"\nCleaned\nSonarr: "+os.Getenv("sonarr_eventtype")+"\nRadarr: "+os.Getenv("radarr_eventtype")+"\n______________\n"))
}

func flagString(fs *string) string {
	return fmt.Sprint(*fs)
}

func flagInt(fi *int64) int64 {
	return int64(*fi)
}

func flagBool(fb *bool) bool {
	return bool(*fb)
}

func writeToFile(path string, data []byte) {
	if filepath.Dir(path) == filepath.Dir(f.Log) {
		// do natta thing
	} else {
		path = filepath.Join(f.Log, filepath.Base(path))
	}
	fmt.Println(path)

	fileName := filepath.Base(path)
	fileName = strings.Replace(fileName, "\"", "", -1)
	fileName = strings.Replace(fileName, ":", "", -1)
	fileName = strings.Replace(fileName, "*", "", -1)
	fileName = strings.Replace(fileName, "?", "", -1)
	fileName = strings.Replace(fileName, "<", "", -1)
	fileName = strings.Replace(fileName, ">", "", -1)
	fileName = strings.Replace(fileName, "|", "", -1)
	fileName = strings.Trim(fileName, " ")

	filePath := filepath.Dir(path)
	path = filepath.Join(filePath, fileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	check(err)

	defer f.Close()

	_, err = f.Write(data) // _ was n
	check(err)
	//fmt.Printf("wrote %d bytes\n", n)

	f.Sync()
}

// Check err
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Only print debug output if the debug flag is true
func printDebug(format string, vars ...interface{}) {
	if *flagDebug {
		if vars[0] == nil {
			fmt.Println(format)
			return
		}
		fmt.Printf(format, vars...)
	}
}

// Hold flag data
type fFlags struct {
	Log   string
	Debug bool
}

// Print the logo, obviously
func printLogo() {
	if *flagNoLogo {
		return
	}
	fmt.Println(" ██████╗██╗     ███╗   ██╗██████╗")
	fmt.Println("██╔════╝██║     ████╗  ██║██╔══██╗")
	fmt.Println("██║     ██║     ██╔██╗ ██║██║  ██║")
	fmt.Println("██║     ██║     ██║╚██╗██║██║  ██║")
	fmt.Println("╚██████╗███████╗██║ ╚████║██████╔╝")
	fmt.Println(" ╚═════╝╚══════╝╚═╝  ╚═══╝╚═════╝ Cleaned")
	fmt.Println("")
}

func buildBody(content map[int]keyVal) (body string) {
	i := 0
	for _ = range content {
		k := content[i].Key
		v := content[i].Value
		body = body + k + ":\t" + v + "\n"
		i++
	}

	return body
}

type keyVal struct {
	Key   string
	Value string
}

func radarrOnDownload() {

	t := time.Now()
	m := make(map[int]keyVal)

	m[0] = keyVal{Key: "ID", Value: os.Getenv("radarr_Movie_Id")}
	title := os.Getenv("radarr_Movie_Title")
	m[1] = keyVal{Key: "Title", Value: title}
	m[2] = keyVal{Key: "Path", Value: os.Getenv("radarr_Movie_Path")}
	m[3] = keyVal{Key: "IMDB", Value: os.Getenv("radarr_Movie_ImdbId")}
	m[4] = keyVal{Key: "File ID", Value: os.Getenv("radarr_MovieFile_Id")}
	m[5] = keyVal{Key: "File Relative Path", Value: os.Getenv("radarr_MovieFile_RelativePath")}
	filePath := os.Getenv("radarr_MovieFile_Path")
	m[6] = keyVal{Key: "File Path", Value: filePath}
	m[7] = keyVal{Key: "File Quality", Value: os.Getenv("radarr_MovieFile_Quality")}
	m[8] = keyVal{Key: "Quality Version", Value: os.Getenv("radarr_MovieFile_QualityVersion")} //1 is the default, 2 for proper, 3+ could be used for anime versions
	m[9] = keyVal{Key: "Release Group", Value: os.Getenv("radarr_MovieFile_ReleaseGroup")}
	m[10] = keyVal{Key: "Scene Name", Value: os.Getenv("radarr_MovieFile_SceneName")}
	sourcePath := os.Getenv("radarr_MovieFile_SourcePath")
	m[11] = keyVal{Key: "Source Path", Value: sourcePath}
	m[12] = keyVal{Key: "Source Folder", Value: os.Getenv("radarr_MovieFile_SourceFolder")}
	m[13] = keyVal{Key: "Time", Value: t.Format("2006-01-02:15:04:05")}

	// Check if Source Path still exists
	if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
		m[14] = keyVal{Key: "Source Exists", Value: "true"}
	} else {
		m[14] = keyVal{Key: "Source Exists", Value: "false"}
	}

	// Check if Destination Path exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		m[15] = keyVal{Key: "Destination Exists", Value: "true"}
	} else {
		m[15] = keyVal{Key: "Destination Exists", Value: "false"}
	}

	// Check if Destination Path exists
	postProcessedPath := strings.Replace(filePath, ".avi", ".mp4", -1)
	postProcessedPath = strings.Replace(postProcessedPath, ".divx", ".mp4", -1)
	postProcessedPath = strings.Replace(postProcessedPath, ".m4v", ".mp4", -1)
	postProcessedPath = strings.Replace(postProcessedPath, ".mkv", ".mp4", -1)

	m[16] = keyVal{Key: "Post-Processed Path", Value: postProcessedPath}
	postProcessedPathExists := false
	if _, err := os.Stat(postProcessedPath); !os.IsNotExist(err) {
		m[17] = keyVal{Key: "Post-Processed Path Exists", Value: "true"}
		postProcessedPathExists = true
	} else {
		m[17] = keyVal{Key: "Post-Processed Path Exists", Value: "false"}
	}

	if postProcessedPathExists {
		err := os.Remove(sourcePath)
		if err != nil {
			m[18] = keyVal{Key: "Source Deleted", Value: "false"}
			m[19] = keyVal{Key: "Source Delete Error", Value: err.Error()}
		} else {
			m[18] = keyVal{Key: "Source Deleted", Value: "true"}
			m[19] = keyVal{Key: "Source Delete Error", Value: "na"}
		}
	}

	fileName := "[Radarr Cleaned] " + title + ".log"
	logPath := filepath.Join(f.Log, fileName)
	writeToFile(logPath, []byte(buildBody(m)))
}

func sonarrOnDownload() {

	t := time.Now()
	m := make(map[int]keyVal)

	m[0] = keyVal{Key: "IsUpgrade", Value: os.Getenv("sonarr_IsUpgrade")} //True when an an existing file is upgraded, otherwise False
	m[1] = keyVal{Key: "ID", Value: os.Getenv("sonarr_Series_Id")}
	title := os.Getenv("sonarr_series_title")
	m[2] = keyVal{Key: "Title", Value: title}
	m[3] = keyVal{Key: "Path", Value: os.Getenv("sonarr_Series_Path")}
	m[4] = keyVal{Key: "TvdbId", Value: os.Getenv("sonarr_Series_TvdbId")}
	m[5] = keyVal{Key: "TvMazeId", Value: os.Getenv("sonarr_Series_TvMazeId")}
	m[6] = keyVal{Key: "IMDB", Value: os.Getenv("sonarr_Series_Imdb")}
	m[7] = keyVal{Key: "Series Type", Value: os.Getenv("sonarr_Series_Type")}
	m[8] = keyVal{Key: "Episode File ID", Value: os.Getenv("sonarr_EpisodeFile_Id")}
	m[9] = keyVal{Key: "Relative Path", Value: os.Getenv("sonarr_EpisodeFile_RelativePath")}
	filePath := os.Getenv("sonarr_EpisodeFile_Path")
	m[10] = keyVal{Key: "File Path", Value: filePath}
	m[11] = keyVal{Key: "Episode Count", Value: os.Getenv("sonarr_EpisodeFile_EpisodeCount")}
	seasonNumber := os.Getenv("sonarr_EpisodeFile_SeasonNumber")
	m[12] = keyVal{Key: "Season Number", Value: seasonNumber}
	episodeNumber := os.Getenv("sonarr_EpisodeFile_EpisodeNumbers")
	m[13] = keyVal{Key: "Episode Number", Value: episodeNumber}
	m[14] = keyVal{Key: "Episode Air Dates", Value: os.Getenv("sonarr_EpisodeFile_EpisodeAirDates")}
	m[15] = keyVal{Key: "Episode Air Dates UTC", Value: os.Getenv("sonarr_EpisodeFile_EpisodeAirDatesUtc")}
	episodeTitle := os.Getenv("sonarr_EpisodeFile_EpisodeTitles")
	m[16] = keyVal{Key: "Episode Title", Value: episodeTitle}
	m[17] = keyVal{Key: "File Quality", Value: os.Getenv("sonarr_EpisodeFile_Quality")}
	m[18] = keyVal{Key: "File Quality Version", Value: os.Getenv("sonarr_EpisodeFile_QualityVersion")} // 1 is the default, 2 for proper, 3+ could be used for anime versions
	m[19] = keyVal{Key: "Release Group", Value: os.Getenv("sonarr_EpisodeFile_ReleaseGroup")}
	sceneName := os.Getenv("sonarr_EpisodeFile_SceneName")
	m[20] = keyVal{Key: "Scene Name", Value: sceneName}
	sourcePath := os.Getenv("sonarr_EpisodeFile_SourcePath")
	m[21] = keyVal{Key: "Source Path", Value: sourcePath}
	m[22] = keyVal{Key: "Source Folder", Value: os.Getenv("sonarr_EpisodeFile_SourceFolder")}
	m[23] = keyVal{Key: "Deleted Relative Paths", Value: os.Getenv("sonarr_DeletedRelativePaths")}
	m[24] = keyVal{Key: "Deleted Paths", Value: os.Getenv("sonarr_DeletedPaths")}
	m[25] = keyVal{Key: "Time", Value: t.Format("2006-01-02:15:04:05")}

	// Check if Source Path still exists
	if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
		m[26] = keyVal{Key: "Source Exists", Value: "true"}
	} else {
		m[26] = keyVal{Key: "Source Exists", Value: "false"}
	}

	// Check if Destination Path exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		m[27] = keyVal{Key: "Destination Exists", Value: "true"}
	} else {
		m[27] = keyVal{Key: "Destination Exists", Value: "false"}
	}

	// Check if Destination Path exists
	postProcessedPath := strings.Replace(filePath, ".avi", ".mp4", -1)
	postProcessedPath = strings.Replace(postProcessedPath, ".divx", ".mp4", -1)
	postProcessedPath = strings.Replace(postProcessedPath, ".m4v", ".mp4", -1)
	postProcessedPath = strings.Replace(postProcessedPath, ".mkv", ".mp4", -1)

	m[28] = keyVal{Key: "Post-Processed Path", Value: postProcessedPath}
	postProcessedPathExists := false
	if _, err := os.Stat(postProcessedPath); !os.IsNotExist(err) {
		m[29] = keyVal{Key: "Post-Processed Path Exists", Value: "true"}
		postProcessedPathExists = true
	} else {
		m[29] = keyVal{Key: "Post-Processed Path Exists", Value: "false"}
	}

	if postProcessedPathExists {
		err := os.Remove(sourcePath)
		if err != nil {
			m[18] = keyVal{Key: "Source Deleted", Value: "false"}
			m[19] = keyVal{Key: "Source Delete Error", Value: err.Error()}
		} else {
			m[18] = keyVal{Key: "Source Deleted", Value: "true"}
			m[19] = keyVal{Key: "Source Delete Error", Value: "na"}
		}
	}

	fileTitle := ""
	if sceneName != "" {
		fileTitle = title + " - " + "S" + seasonNumber + "E" + episodeNumber + " - " + sceneName
	} else {
		fileTitle = title + " - " + "S" + seasonNumber + "E" + episodeNumber + " - " + episodeTitle
	}
	fileName := "[Sonarr Cleaned] " + fileTitle + ".log"
	logPath := filepath.Join(f.Log, fileName)
	writeToFile(logPath, []byte(buildBody(m)))
}
