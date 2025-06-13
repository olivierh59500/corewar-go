;redcode
;name Aggressive
;author Core War Community
;strategy Very aggressive bombing pattern
;assert 1

        SPL 4            ; Create multiple processes
        SPL 2            ; More processes
        JMP imp1         ; Some become imps
        JMP imp2         ; Others bomb
imp1:   MOV 0, 1         ; Imp that goes forward
imp2:   MOV 0, -1        ; Imp that goes backward

end