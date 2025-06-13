;redcode
;name ImpStomper
;author Core War Community
;strategy Moves through memory laying DAT bombs
;assert 1

        MOV 2, >2        ; Copy bomb forward using post-increment
        JMP -1           ; Loop
ptr:    DAT #0, #1       ; Bomb and pointer

end