7
(defun add
  (lambda (a)
    (lambda (b)
      (+ a b))))
((add 3) 4)