24
(defun fact
  (lambda (n)
    (cond
      ((eq n 0) 1)
      (#t (* n (fact (- n 1)))))))
(defun Y
  (lambda (h)
    (lambda (f)
      (h (lambda v ((f f) v) )) )
    (lambda (f)
      (h (lambda v ((f f) v) )) ) ))
(fact 4)
