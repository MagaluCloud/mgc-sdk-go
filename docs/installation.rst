Installation
============

Requirements
------------

- Go 1.21 or higher
- A Magalu Cloud account with API access

Installation
-----------

Install the MGC Go SDK using Go modules:

.. code-block:: bash

   go get github.com/MagaluCloud/mgc-sdk-go

Or add it to your `go.mod` file:

.. code-block:: go

   require github.com/MagaluCloud/mgc-sdk-go v0.1.0

Then run:

.. code-block:: bash

   go mod tidy

Verify Installation
------------------

Create a simple test to verify the installation:

.. code-block:: go

   package main

   import (
       "fmt"
       "github.com/MagaluCloud/mgc-sdk-go/client"
   )

   func main() {
       fmt.Println("MGC Go SDK installed successfully!")
       
       // Test client creation
       c := client.NewMgcClient("test-token")
       fmt.Printf("Client created: %v\n", c != nil)
   }

Run the test:

.. code-block:: bash

   go run main.go

You should see:

.. code-block:: text

   MGC Go SDK installed successfully!
   Client created: true

Next Steps
----------

After installation, you can:

1. Set up authentication (see :doc:`authentication`)
2. Configure the client (see :doc:`configuration`)
3. Start using the services (see :doc:`services/index`)
4. Check out examples (see :doc:`examples`) 