12
(defun mul
  (lambda (a b)
      (cond
        ((eq b 0) 0)
        (#t (+ a (mul a (- b 1)))))))
(mul 3 4)
