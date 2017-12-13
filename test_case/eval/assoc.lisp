a

(defun caar (lst)
  (car (car lst)))
(defun cadar (lst)
  (car (cdr (car lst))))
(defun assoc. (x y)
  (cond ((eq (caar y) x) (cadar y))
        (#t (assoc. x (cdr y)))))

(assoc. 'x ('((x a) (y b))))
