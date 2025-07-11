Authentication
==============

The MGC Go SDK uses API tokens for authentication. This section explains how to obtain and use API tokens with the SDK.

Obtaining an API Token
---------------------

To use the MGC Go SDK, you need an API token from your Magalu Cloud account:

1. Log in to the `Magalu Cloud Console <https://console.magalu.cloud>`_
2. Navigate to **Settings** > **API Keys**
3. Click **Create API Key**
4. Give your key a name and description
5. Copy the generated token (you won't be able to see it again)

For more detailed information about API keys, see the `Magalu Cloud documentation <https://docs.magalu.cloud/docs/devops-tools/api-keys/overview>`_.

Using API Tokens
---------------

Basic Authentication
~~~~~~~~~~~~~~~~~~~

The simplest way to authenticate is to pass your API token directly to the client:

.. code-block:: go

   package main

   import (
       "github.com/MagaluCloud/mgc-sdk-go/client"
   )

   func main() {
       apiToken := "your-api-token-here"
       c := client.NewMgcClient(apiToken)
       
       // Use the client...
   }

Environment Variables
~~~~~~~~~~~~~~~~~~~~

For security, it's recommended to store your API token in an environment variable:

.. code-block:: go

   package main

   import (
       "os"
       "github.com/MagaluCloud/mgc-sdk-go/client"
   )

   func main() {
       apiToken := os.Getenv("MGC_API_TOKEN")
       if apiToken == "" {
           panic("MGC_API_TOKEN environment variable is required")
       }
       
       c := client.NewMgcClient(apiToken)
       
       // Use the client...
   }

Set the environment variable:

.. code-block:: bash

   export MGC_API_TOKEN="your-api-token-here"

Or for Windows:

.. code-block:: cmd

   set MGC_API_TOKEN=your-api-token-here

Configuration Files
~~~~~~~~~~~~~~~~~~

You can also load the API token from a configuration file:

.. code-block:: go

   package main

   import (
       "os"
       "bufio"
       "strings"
       "github.com/MagaluCloud/mgc-sdk-go/client"
   )

   func loadConfig(filename string) (map[string]string, error) {
       config := make(map[string]string)
       file, err := os.Open(filename)
       if err != nil {
           return nil, err
       }
       defer file.Close()

       scanner := bufio.NewScanner(file)
       for scanner.Scan() {
           line := strings.TrimSpace(scanner.Text())
           if line != "" && !strings.HasPrefix(line, "#") {
               parts := strings.SplitN(line, "=", 2)
               if len(parts) == 2 {
                   config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
               }
           }
       }
       return config, scanner.Err()
   }

   func main() {
       config, err := loadConfig(".env")
       if err != nil {
           panic(err)
       }
       
       apiToken := config["MGC_API_TOKEN"]
       c := client.NewMgcClient(apiToken)
       
       // Use the client...
   }

Token Security Best Practices
---------------------------

1. **Never commit tokens to version control**
   - Use environment variables or secure configuration management
   - Add `.env` files to `.gitignore`

2. **Use different tokens for different environments**
   - Development tokens for testing
   - Production tokens for live applications

3. **Rotate tokens regularly**
   - Delete old tokens when they're no longer needed
   - Create new tokens periodically

4. **Limit token permissions**
   - Only grant the minimum required permissions
   - Use read-only tokens when possible

5. **Monitor token usage**
   - Check the Magalu Cloud console for token usage
   - Set up alerts for unusual activity

Example: Secure Token Management
------------------------------

Here's a complete example showing secure token management:

.. code-block:: go

   package main

   import (
       "context"
       "log"
       "os"
       "time"

       "github.com/MagaluCloud/mgc-sdk-go/client"
       "github.com/MagaluCloud/mgc-sdk-go/compute"
   )

   func main() {
       // Load API token from environment
       apiToken := os.Getenv("MGC_API_TOKEN")
       if apiToken == "" {
           log.Fatal("MGC_API_TOKEN environment variable is required")
       }

       // Create client with timeout and retry configuration
       c := client.NewMgcClient(
           apiToken,
           client.WithTimeout(30*time.Second),
           client.WithRetryConfig(3, 1*time.Second, 30*time.Second, 1.5),
       )

       // Test the connection
       computeClient := compute.New(c)
       ctx := context.Background()
       
       _, err := computeClient.Instances().List(ctx, compute.ListOptions{})
       if err != nil {
           log.Fatalf("Failed to connect to Magalu Cloud: %v", err)
       }

       log.Println("Successfully authenticated with Magalu Cloud")
   }

Next Steps
----------

After setting up authentication:

1. Configure the client (see :doc:`configuration`)
2. Choose a region (see :doc:`configuration#regions`)
3. Start using the services (see :doc:`services/index`) 