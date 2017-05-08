package main

//// BuildForever watches the filesystem and builds a new binary if something changes.
//// It notifies a channel that a build was created
//func BuildForever(builds chan<- bool) {
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	defer watcher.Close()
//
//	go func() {
//		for {
//			select {
//			case event := <-watcher.Events:
//				if event.Op != fsnotify.Chmod && event.Name != "" {
//					start := time.Now()
//					if err := build(); err == nil { // only notify and log if binary was created successfully.
//						log.Println("built a new binary in", time.Since(start))
//						builds <- true // notify that a new build was successfully created
//					}
//					watcher.Remove(event.Name)
//					watcher.Add(event.Name)
//				}
//			case err := <-watcher.Errors:
//				log.Println(err)
//			}
//		}
//	}()
//
//	files, err := findGoFiles()
//	if err != nil {
//		log.Println(err)
//	}
//
//	for _, file := range files {
//		if err := watcher.Add(filepath.Join(".", file)); err != nil {
//			log.Println(err)
//		}
//	}
//
//	select {}
//}
//
//func findGoFiles() ([]string, error) {
//	var files []string
//	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
//		if strings.HasPrefix(path, "cmd/stratus") { // don't watch stratus itself
//			return nil
//		}
//		if strings.HasSuffix(path, ".go") {
//			files = append(files, path)
//		}
//		return nil
//	})
//	return files, err
//}
