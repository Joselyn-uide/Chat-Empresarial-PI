# Proyecto Integrador: Sistema de chat empresarial Cliente-Servidor con red segmentada

## Integrantes
* Joselyn Paulet Laguna Puetate 
* Jose Emiliano Velasco Sosa 
* Darwin Oswaldo Román Aguirre 

---

## Descripción del problema
Actualmente la empresa presenta dificultades en la comunicación interna entre los colaboradores que trabajan en sus dos sedes. Gran parte de la información operativa se transmite mediante llamadas, mensajes personales o correos electrónicos no centralizados, lo que provoca retrasos, pérdida de información y dificultad para mantener un registro ordenado de las conversaciones.

## Justificación del sistema seleccionado
Se desarrolló una aplicación de chat corporativo en arquitectura cliente-servidor para permitir a los empleados de ambas sedes comunicarse de forma directa, centralizada y sin depender de aplicaciones personales. El sistema utiliza el lenguaje Go por su alta eficiencia en el manejo de concurrencia mediante sockets y goroutines, y se integra con un microservicio de autenticación en FastAPI para proteger el acceso mediante hashes seguros.

## Objetivos del Proyecto
* **Objetivo General:** Desarrollar una aplicación de chat corporativo en arquitectura cliente/servidor, en el lenguaje Go, que permita la interacción concurrente y segura entre usuarios de dos unidades organizativas distintas, con verificación de identidad y sobre una infraestructura de red segmentada.
* **Objetivos Específicos:** 
  1. Implementar funciones de mensajería bidireccional y gestión del historial básico.
  2. Integrar el servicio de autenticación con FastAPI.
  3. Controlar el acceso mediante hashing seguro.
  4. Organizar el código bajo principios de Programación Orientada a Objetos (POO).
  5. Validar la conectividad mediante una red segmentada por VLANs y enrutamiento OSPF.

## Alcance del Proyecto
El proyecto abarca el diseño e implementación de un sistema que combina software y redes:
* **Software:** Lógica de chat cliente-servidor en Go y servicio de autenticación en FastAPI.
* **Infraestructura de Red:** Simulación en Cisco Packet Tracer de dos sedes interconectadas mediante VLANs, enlaces troncales (Trunk), Router-on-a-Stick, enrutamiento dinámico OSPF y redundancia RSTP.

## Exclusiones del Sistema
* No incluye despliegue en nubes comerciales de pago (AWS, Azure, Google Cloud); validación estrictamente local.
* No incluye llamadas de voz, videollamadas ni transmisión de archivos pesados (limitado a texto).
* No incluye aplicaciones móviles nativas ni pasarelas de mensajería pública externa.

---

## Video
Video grupa: https://drive.google.com/file/d/1BgLwMJoLMNJwATurPaoKMaBrOg6PJpEg/view?usp=drivesdk
