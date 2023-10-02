resource "kubernetes_namespace" "app-go" {
  metadata {
    name = "app-go"
  }
}


resource "kubernetes_namespace" "monitoring" {
  metadata {
    name = "monitoring"
  }
}

resource "helm_release" "prometheus" {
  name = "prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart = "kube-prometheus-stack"
  namespace = kubernetes_namespace.monitoring.metadata[0].name
  version = "48.0.0"
  force_update = true 


  set {
    name = "prometheus.service.type"
    value = "LoadBalancer"
  }

  set {
    name = "grafana.service.type"
    value = "LoadBalancer"
  }
    
  set {
    name = "grafana.adminPassword"
    value = "admin"
  }

  set {
    name = "prometheus.prometheusSpec.serviceMonitorNamespaceSelector"
    value = false
  }

  set {
    name = "prometheus.prometheusSpec.serviceMonitorSelector"
    value = false
  }

  set {
    name = "prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues"
    value = false
  }

  set {
    name = "prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues"
    value = false
  }

}

locals {
  api_manifests = {
    "deployment" = file("${path.module}/manifests/api/deployment.yml"), 
    "service" = file("${path.module}/manifests/api/service.yml"),
    "service-monitor" = file("${path.module}/manifests/api/service-monitor.yml"),
  }
}

resource "kubectl_manifest" "app-go" {
  for_each = local.api_manifests
  
  yaml_body = each.value

  depends_on = [
    kubernetes_namespace.app-go,
  ]

}

resource "kubectl_manifest" "prometheus" {
  yaml_body = file("${path.module}/manifests/prometheus.yml")
  
  depends_on = [
    helm_release.prometheus
  ]

}


