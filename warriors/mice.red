;redcode
;name Mice
;author Chip Wendell
;strategy Replicator that spawns many processes
;assert 1

start:  SPL 2           ; Split into two processes
        JMP -1          ; Loop
        MOV 0, 1        ; Copy forward

end start