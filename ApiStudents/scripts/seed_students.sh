#!/usr/bin/env bash
# Script simples para popular a API ApiStudents com dados de exemplo
# Requisitos: servidor ApiStudents rodando em http://localhost:8080
set -euo pipefail

echo "Criando estudantes de exemplo..."

curl -s -X POST http://localhost:8080/students -H 'Content-Type: application/json' -d '{"Name":"João Silva","CPF":123456789,"Email":"joao@example.com","Age":21,"Active":true}'
echo
curl -s -X POST http://localhost:8080/students -H 'Content-Type: application/json' -d '{"Name":"Maria Lima","CPF":987654321,"Email":"maria@example.com","Age":22,"Active":true}'
echo
curl -s -X POST http://localhost:8080/students -H 'Content-Type: application/json' -d '{"Name":"Carlos Souza","CPF":111222333,"Email":"carlos@example.com","Age":20,"Active":false}'
echo

echo "Estudantes criados. Verifique com: curl http://localhost:8080/students"
