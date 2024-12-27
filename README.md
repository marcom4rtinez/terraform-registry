Terraform Registry
---

Go Application to self-host a Terraform Registry, with the possibility to add providers using the API endpoint `/v1/providers/:namespace/:name/upload`


Run the Webserver using the following command:
```bash
go run github.com/marcom4rtinez/terraform-registry@latest
```

Defaults to Port 8080.



> [!WARNING]
> This Application requires no authorization for adding Providers. This is not a production-ready Application