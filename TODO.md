RESTful Interface
=================

Changes to Backends
-------------------

* All backends need a Configure method to apply a configuration. This
  should restart services as needed

* Sig Backend: needs a connect/reconnect method

* VPN Backend: No major changes aside from config

* Stun Backend: No major changes aside from config

Changes to TVPN
---------------

* Needs a Configure method. This should also configure backends and test
  changes in friendship and handle connections accordingly

* Needs an HTTP listener that accepts reconfiguration commands and
  commands to query the internal state. All data transfer will be
  encoded in JSON

Web Interface
=============

This will lean heavily on the RESTful interface. More details to be
added when that is complete.

