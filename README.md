# dps_overlay
A simple damage prediction overlay for League of Legends

## Upcoming features
- accounting for lethality/pen in ability damage predictions
- auto-update for JSON file data for each new patch
- item recommendation based on damage
- dps calculation (attack speed/combos)
- calibration for niche abilities (like passives) that don't work as intended
## Used libraries
- Go bindings for raylib: https://github.com/gen2brain/raylib-go
- syscall package for interacting with Windows
- lolstaticdata, a library that auto-generates JSON files with game data using https://wiki.leagueoflegends.com: https://github.com/meraki-analytics/lolstaticdata
## Required settings
### For the overlay to work:
- in-game settings->Video->General->Window Mode->Borderless
### To make the overlay transparent instead of black, set either:
1) Settings->System->Display->Graphics->Add Desktop App (browse the app)->GPU preference->set to a non-nvidia gpu
2) Nvidia Control Panel->Manage 3D settings->OpenGL GDI compatibility->Prefer compatible

dps_overlay is not endorsed by Riot Games and does not reflect the views or opinions of Riot Games or anyone officially involved in producing or managing Riot Games properties. Riot Games and all associated properties are trademarks or registered trademarks of Riot Games, Inc