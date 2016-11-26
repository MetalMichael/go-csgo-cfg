Based Heavily on [go-ini](https://github.com/go-ini/ini)

Package ini provides INI file read and write functionality in Go.

## Example Config format

    ammo_grenade_limit_default 1
    ammo_grenade_limit_flashbang 2
    ammo_grenade_limit_total 4
    bot_quota "0"                 // Determines the total number of bots in the game
    cash_player_bomb_defused 300
    cash_player_bomb_planted 300
    ammo_grenade_limit_total 5

## Removed Features

- As CSGO doesn't support them, there is no longer section support. Could potentially be added back in, but has been pulled out for now.
- No longer supports embedded or nested structs, again not needed for this purpose
- Booleans
- Arrays
- Auto Increment


## Installation

To use with latest changes:

	go get github.com/metalmichael/go-csgo-cfg

Please add `-u` flag to update in the future.

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
