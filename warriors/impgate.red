;redcode
;name ImpGate
;author Core War Community
;strategy Creates a proper imp gate with multiple processes
;assert 1

        SPL 2            ; Split into multiple processes
        JMP -1           ; Some processes loop here
        MOV bomb, 1      ; Others create the gate
        JMP -1           ; And maintain it
bomb:   DAT #0, #0       ; The bomb

end