(ns expose.main
  (:gen-class))

(defn -main [& args]
  (prn "hi"))

(defn f [x y]
  (* x y))
