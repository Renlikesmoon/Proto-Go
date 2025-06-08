module github.com/Renlikesmoon/Proto-Go // Your module path

go 1.21 // Go version

require (
        go.mau.fi/whatsmeow v0.0.0-20250606170101-3afe34f8ab8f // Dependency
        github.com/mdp/qrterminal/v3 v3.2.1 // Dependency
)

replace go.mau.fi/whatsmeow => github.com/tulir/whatsmeow v0.0.0-20250606170101-3afe34f8ab8f // Replacement directive