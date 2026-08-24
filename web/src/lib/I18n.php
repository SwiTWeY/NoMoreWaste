<?php

class I18n
{
    private static array $t = [];
    private static string $langue = 'fr';

    public static function init(): void
    {
        if (isset($_GET['lang']) && in_array($_GET['lang'], ['fr', 'en'], true)) {
            $_SESSION['lang'] = $_GET['lang'];
        }
        self::$langue = $_SESSION['lang'] ?? 'fr';
        self::$t = require __DIR__ . '/../lang/' . self::$langue . '.php';
    }

    public static function langue(): string
    {
        return self::$langue;
    }

    public static function t(string $cle): string
    {
        return self::$t[$cle] ?? $cle;
    }
}

function t(string $cle): string
{
    return I18n::t($cle);
}
