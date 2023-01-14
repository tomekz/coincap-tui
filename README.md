## coincap-tui
coincap-tui let's you check crypto prices in your terminal.

Features:
- fetch crypto assets data from [ coincap ](https://docs.coincap.ioA/) REST API
- display results in tabular format
- nice UI with [bubble-tea](https://github.com/charmbracelet/bubbletea)

<img src="img/table" alt="demo" />

table with crypto assets

<img src="img/graph" alt="demo" />

price history for the last 14 days

## :keyboard: keybindings

|      Key      |                Description                |
| :-----------: | :---------------------------------------: |
|     `r`       |           refresh data                    |
|   `enter`     |           show price history graph.       |
|      `b`      |           go back                         |
|     `j`       |             go down                       |
|     `k`       |              go up                        |
| `g, home`     |         go to top                         |
| `G, end`      |        go to end                          |
| `ctrl-c`      |                exit                       |


## how to run

```sh
 go run main.go
```

## acknowledgments

Inspired by [tinance](https://github.com/Alcadramin/tinance)
