#!/bin/bash

generatePdf () {
    cd $1
    pandoc Docu.md -o Dokumentation.pdf --from markdown --template eisvogel --listings
    cd ..
}

generatePdf "./datalake" 
generatePdf "./github-webhook-lambda-sns" 
generatePdf "./fargate" 