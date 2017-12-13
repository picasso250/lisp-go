3
(defun len
  (lambda (lst)
    (cond ((eq lst ()) 0)
          (#t (+ 1 (len (cdr lst)))))))
(len ('(a b c)))