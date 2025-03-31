(ns clojure-noob.core
  (:gen-class))

(def map-in-a-map {:outer-key "outer value" :inner-map {:inner-key "inner value"}})
(def mixed-vector ["string preceded by ratio" 1/4 "another string then floating point" 0.25])

(defn -main
  "???"
  [& args]
  ((if true
     (do (println "First thing in do")
         (println "second thing in do")
         "Returned from do")
     (println "Else branch"))) ; if must have a else branch
  (when true
    (println "First true thing in when")
    (println "Second true thing in when")
    "Returned from when if condition is true, else nil") ; when can be used when there is no else branch and nil is returned for falsy
  (when (nil? nil) "nil? is a function that returns true if its argument is nil, false otherwise")
  (if (= 1 2) ; (= arg1 arg2) is a function that returns true if arg1 is equal to arg2, false otherwise
    (println "1 is never equal to 2")
    (println "indeed 1 is not equal to 2"))
  (or "or returns this string, the first truthy value"
      "not the second"
      "nor the third and for sure not the forth falsy value" nil)
  (if (or false nil) "" "or with only falsy values returns the last falsy value")
  (println (and "first truthy value" "and returns this last truthy value"))
  (if (and true false) "" "and with mixed truthy and falsy values returns the first falsy value")
  (println (get map-in-a-map :inner-map))    ; get and get-in are functions that return the value of a key in a map
  (println (get-in map-in-a-map [:inner-map :inner-key])) ; get-in can access nested maps
  (println (:inner-key map-in-a-map))        ; :keyword is a shorthand for (get map :key)
  (println (get mixed-vector 1))             ; get also used for vectors
  (println (get 0 (vector (vector "inner vector with 3 item" 1 nil) "outer vector with 2 items" "second item")))
  (println (conj mixed-vector "added to the end of the vector")) ; conj is a function that adds an item to the end of a collection
  (println (list "list of" 3 "items"))       ; list is a function that creates a list
  (println '(1 2 3))                         ; ' is a shorthand for list
  (println (when (= (nth '(1 2 3) 1) 2) "item at index 1 in the list is indeed")) ; nth is a inefficient way to access list by index
  (println (conj '(1 2 3) 0))                ; conj can also add an item to the beginning of a list; what?
  (println (if (= #{true "2" 3} (hash-set true "2" 3)) "literal set equal with hash set" "")) ; #{} is a set literal
  (when (contains? #{true "2" 3} "2") "found it") ; contains? returns boolean when checking for a value
  (println (get #{true "2" 3} "2"))          ; get can be used to get a value from a set
  (:keyword #{:keyword "value"})           ; :keyword can be used to get a value from a set
  )


(def is-ratio-equal-with-float (= 1/4 0.25))

(defn function-name
  "optional docstring"
  ([parameter1 parameter2]
   (if parameter1
     (println "parameter1 is truthy")
     (println parameter2)))
  ([parameter1]
   (function-name parameter1 "default parameter2")))

(defn function-with-varargs
  ([parameter1 & rest]
   (println "parameter1:" parameter1)
   (println "rest:" rest)))