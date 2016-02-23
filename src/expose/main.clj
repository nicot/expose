(ns expose.main
  (:gen-class))

;; https://systembash.com/a-simple-java-tcp-server-and-tcp-client/
(require '[clojure.java.io :as io])

(import '[java.net ServerSocket])

(defn expose
  "Expose a service to the world!"
  [host lport & args]
  (prn lport))

(defn receive
  "read a line of text from the socket"
  [socket]
  (.readLine (io/reader socket)))

(defn send
  "send the message out the socket"
  [socket msg]
  (let [writer (io/writer socket)]
    (.write writer msg)
    (.flush writer)))

(defn handler
  "simple echo handler"
  [msg]
  msg)

(defn serve
  "Allow exposing ports on this host"
  [port & args]
  (with-open [server-sock (ServerSocket. port)
              sock (.accept server-sock)]
    (let [msg-in (receive sock)
          msg-out (handler msg-in)]
      (send sock msg-out))))

(defn -main [& args]
  (case (first args)
    "serve" (serve args)
    (expose args)))
