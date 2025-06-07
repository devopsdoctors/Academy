# 📅 Course Schedule - Sepahsalar Group

This course is designed for individuals who are already familiar with basic DevOps concepts, and our goal is to dive deeper into advanced topics and gain a deeper understanding of various technologies.  
Below is the schedule for the first two weeks of the Sepahsalar group course.

---

## 🌐 Module 1: Core Networking Concepts & Web Architecture
### **Course Topics**
1. **Introduction to Web Servers, Load Balancers, and API Gateways**
   - Why we need web servers
   - When to use load balancers
   - When to use API gateways
   - Hardware vs. software load balancers

2. **HTTP Protocol**
   - HTTP basics
   - HTTP/1.1 vs HTTP/2 vs HTTP/3

3. **TLS and SSL**
   - How TLS works
   - Differences between TLS versions

4. **Networking Fundamentals**
   - OSI vs TCP/IP Model
   - Proxy vs NAT
   - DNS and DoH (DNS over HTTPS)
   - How DNS works in Linux
   - How CDNs work

### **Training & Hands-on**
1. Deep dive into TLS and its versions
2. How the Linux kernel handles TCP connections
3. Trace an HTTP request across the OSI layers
4. Build a simple web server from scratch
5. Compare tunneling protocols: IPIP, Geneva, VXLAN
6. Understand Linux bridges and their use
7. Introduction to Anycast and BGP
8. What is Software-Defined Networking (SDN)
9. Site-to-Site VPN and IPsec

---
## 🖥️ Module 2: Virtualization & Containers
### **Course Topics**

1. **Introduction to Virtualization**
   - What is virtualization
   - KVM and MicroVM

2. **Introduction to Containers**
   - What are containers
   - Namespaces and cgroups (v1 vs v2)
   - Container internals

3. **Hands-On: Creating Containers with Linux**
   - Using `unshare` and Linux namespaces
   - Basic container shell without Docker

4. **Docker Concepts**
   - Docker overview
   - Docker images and layering
   - Container runtimes (runc, containerd, etc.)
   - Docker container lifecycle
   - Docker storage types
   - Docker networking types

### **Training & Hands-on**
1. What is Copy-On-Write (COW) in Docker storage
2. Create a Docker container manually (without Docker CLI)
3. Manually connect two containers running in different VMs
4. Share namespaces between containers and understand how it works
5. Container security best practices (image signing, user management, seccomp, etc.)
6. Dockerfile and Docker-compose best practices

---

## ☸️ Module 3: Kubernetes Deep Dive
### **Course Topics**

1. **Kubernetes Architecture**
   - Components of the control plane and node
   - etcd, API server, controller manager, scheduler, kubelet, kube-proxy

2. **Workload Management**
   - Pods, ReplicaSets, Deployments
   - StatefulSets vs Deployments
   - DaemonSets
   - Jobs and CronJobs

3. **Kubernetes Networking**
   - Services: ClusterIP, NodePort, LoadBalancer, ExternalName
   - Ingress and Ingress Controllers
   - CNI plugins and K8s networking model


4. **Storage**
   - PersistentVolumes (PV) and PersistentVolumeClaims (PVC)
   - StorageClasses and dynamic provisioning

5. **Security**
   - RBAC (Role-Based Access Control)
   - Network Policies
   - Secrets and ConfigMaps

6. **Packaging and Deployment**
   - Helm and Helm Charts
   - GitOps basics and tools (ArgoCD, Flux)

### **Training & Hands-on**
1. Deploy k8s cluster
2. Deploy wordpress on Kubernetes
3. different rollout strategy on Kubernetes
4. cilium CNI 
4. Calico CNI

---

Good luck and get ready for an exciting and educational journey! 🚀
