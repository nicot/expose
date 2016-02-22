(set-env!
 :resource-paths #{"src"})

(task-options!
 pom {:project 'my-project
      :version "0.1.0"}
 jar {:main 'expose.main})


(require
 '[expose.main])

(deftask build
  "build my project"
  []
  (comp (pom) (jar) (install) (target)))

(deftask run []
  (expose.main/-main))
