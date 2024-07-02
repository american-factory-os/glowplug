# Generate OPC UA SSL certificates

To generate a OPC UA client certificate:

1. Replace `urn:HOSTNAME.mshome.net:OPCUA:SimulationServer` with the correct OPC UA server URI.
2. Run the openssl script:
    ```bash
    ./gen.sh
    ```

This will generate
* `public.der` - public certificate in DER format
* `default_pk.pem` - private key

Don't forget to have your OPC UA server trust this certificate.

