;redcode
;name Sucker
;author Core War Community
;strategy Creates an imp-killing spiral
;assert 1

start:  JMP 2            ; Jump over the gate
        DAT #0, #0       ; The gate
        SPL -1           ; Create processes that go back to the gate
        MOV -2, <-3      ; Build the spiral
        JMP -1           ; Keep building

end start