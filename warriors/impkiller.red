;redcode
;name ImpKiller
;author Core War Community
;strategy Specifically designed to kill imps with a gate
;assert 1

        MOV 3, 4         ; Copy the DAT bomb forward
        MOV 2, 3         ; Copy it again to create a wall
        JMP -2           ; Jump back to continue
        DAT #0, #0       ; The bomb that kills imps

end