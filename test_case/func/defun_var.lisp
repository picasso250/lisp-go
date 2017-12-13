(c d)
(defun cadr
  (lambda (e) (car (cdr e))))
(cadr ('((a b) (c d) e)))