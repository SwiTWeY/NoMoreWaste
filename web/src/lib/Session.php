<?php

class Session
{
    public static function demarrer(): void
    {
        if (session_status() === PHP_SESSION_NONE) {
            session_start();
        }
    }

    public static function connecter(string $token, array $utilisateur): void
    {
        $_SESSION['token'] = $token;
        $_SESSION['utilisateur'] = $utilisateur;
    }

    public static function token(): ?string
    {
        return $_SESSION['token'] ?? null;
    }

    public static function utilisateur(): ?array
    {
        return $_SESSION['utilisateur'] ?? null;
    }

    public static function estConnecte(): bool
    {
        return isset($_SESSION['token']);
    }

    public static function deconnecter(): void
    {
        $_SESSION = [];
        session_destroy();
    }
}
