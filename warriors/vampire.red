;redcode
;name Vampire
;author Core War Community
;strategy Converts enemy processes by making them jump to a pit
;assert 1

pit:    JMP pit         ; Vampire pit - trap for enemy processes

        MOV pit, @fang  ; Place pit at target location
        ADD #10, fang   ; Move to next target
        JMP -2          ; Continue vampiring
fang:   DAT #0, #10     ; Fang pointer

end