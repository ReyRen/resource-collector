#!/bin/bash
kubectl port-forward dcgm-exporter-1608880229-9zm2b --address=172.18.29.80 9401:9400 &
kubectl port-forward dcgm-exporter-1608880229-l5njt --address=172.18.29.80 9402:9400 &
kubectl port-forward dcgm-exporter-1608880229-lxmt5 --address=172.18.29.80 9403:9400 &
kubectl port-forward dcgm-exporter-1608880229-pstc4 --address=172.18.29.80 9404:9400 &
kubectl port-forward dcgm-exporter-1608880229-7mfcv --address=172.18.29.80 9405:9400 &
kubectl port-forward dcgm-exporter-1608880229-s6wgv --address=172.18.29.80 9406:9400 &
kubectl port-forward dcgm-exporter-1608880229-nw26s --address=172.18.29.80 9407:9400 &
kubectl port-forward dcgm-exporter-1608880229-ttw4z --address=172.18.29.80 9408:9400 &