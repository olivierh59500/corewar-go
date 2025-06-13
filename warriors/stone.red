;redcode
;name Stone
;author Core War Community
;strategy Simple stone that bombs forward continuously
;assert 1

        MOV <2, 3        ; Copy and decrement
        ADD #4, 1        ; Increment the pointer
        JMP -2           ; Loop back
        DAT #0, #0       ; Bomb to copy

end