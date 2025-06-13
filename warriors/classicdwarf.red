;redcode
;name ClassicDwarf
;author A.K. Dewdney
;strategy Original 1984 Dwarf that bombs every 4 locations
;assert 1

        ADD #4, 3        ; Add 4 to the DAT pointer
        MOV #0, @2       ; Bomb the location pointed to by DAT
        JMP -2           ; Go back to ADD
        DAT #0, #0       ; This will be the pointer

end