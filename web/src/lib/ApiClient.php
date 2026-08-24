<?php

class ApiClient
{
    private string $baseUrl;
    private ?string $token;

    public function __construct(string $baseUrl, ?string $token = null)
    {
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->token = $token;
    }

    public function get(string $chemin): array
    {
        return $this->requete('GET', $chemin, null);
    }

    public function post(string $chemin, array $corps): array
    {
        return $this->requete('POST', $chemin, $corps);
    }

    public function telecharger(string $chemin): array
    {
        $ch = curl_init($this->baseUrl . $chemin);

        $entetes = [];
        if ($this->token !== null) {
            $entetes[] = 'Authorization: Bearer ' . $this->token;
        }

        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $entetes);

        $corps = curl_exec($ch);
        $statut = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $type = curl_getinfo($ch, CURLINFO_CONTENT_TYPE);
        curl_close($ch);

        return [
            'statut' => $statut,
            'corps' => $corps,
            'type' => $type,
        ];
    }

    private function requete(string $methode, string $chemin, ?array $corps): array
    {
        $ch = curl_init($this->baseUrl . $chemin);

        $entetes = ['Accept: application/json'];
        if ($this->token !== null) {
            $entetes[] = 'Authorization: Bearer ' . $this->token;
        }
        if ($corps !== null) {
            $entetes[] = 'Content-Type: application/json';
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($corps));
        }

        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $methode);
        curl_setopt($ch, CURLOPT_HTTPHEADER, $entetes);

        $reponse = curl_exec($ch);
        $statut = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        return [
            'statut' => $statut,
            'donnees' => json_decode($reponse, true),
        ];
    }
}
