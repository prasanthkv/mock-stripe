# mock-stripe  ![version](https://img.shields.io/badge/version-1.0.1--Mock-orange.svg?style=flat) ![License: MIT](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)


mock-stripe is a mock HTTP server that responds like the real Stripe API. It
can be used instead of Stripe's testmode to make test suites integrating with
Stripe faster and less brittle.

stripe-mock is powered by [GO-Lang](http://www.golangbootcamp.com/book/intro),It operates limited statefulness with 
(i.e. it remember new resources that are created with it for few hours) and responds
with sample data that's generated using a similar scheme to the one found in
the [API reference](https://stripe.com/docs/api).

## Usage

Get it from Homebrew

~~~ sh
brew install go --cross-compile-common
~~~

Or if you have Go installed you can build it:

~~~ sh
go get 

go build
~~~

Run it:

~~~ sh
mockstripe
~~~

Then from another terminal:

~~~ sh
curl -i http://localhost:8080/v1/version -H "Authorization: Bearer sk_test_123"
~~~

By default, stripe-mock runs on port 8080, but is configurable with the
`-port` option. (TODO)

## Development
TBD

## Supported Operations

- [x] Auth
- [x] Capture
- [x] Refund

### Auth
~~~
curl -X POST \
  http://localhost:8080/v1/charges \
  -H 'accept: application/json' \
  -H 'authorization: Bearer sk_test_0zzzz0zXXzOXXXX4X00zXzz0' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'idempotency-key: a_1307' \
  -d 'capture=false&amount=1000&currency=usd&destination%5Baccount%5D=acct_1AsDVNAeE9ZhXsLk&destination%5Bamount%5D=1000&source%5Baddress_line1%5D=2145%20Hamilton%20Avenue&source%5Baddress_city%5D=San%20Jose&source%5Bobject%5D=card&source%5Bnumber%5D=4000000000000077&source%5Bexp_year%5D=2022&source%5Bexp_month%5D=1&source%5Bname%5D=QIB&source%5Baddress_state%5D=CA&source%5Baddress_zip%5D=95125&source%5Baddress_country%5D=US&description=eBay%3A%20pkv_usa'
~~~

### Capture
~~~
curl -X POST \
  http://localhost:8080/v1/charges/ch_1TESTAGQ2G0H1tnT4CMErOEL/capture \
  -H 'accept: application/json' \
  -H 'authorization: Bearer sk_test_0zzzz0zXXzOXXXX4X00zXzz0' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'idempotency-key: c_1308' \
  -d 'amount=500&destination%5Bamount%5D=500'
~~~

### Refund
~~~
curl -X POST \
  http://localhost:8080/v1/charges/ch_1TESTAGQ2G0H1tnT4CMErOEL/refunds \
  -H 'accept: application/json' \
  -H 'authorization: Bearer sk_test_0zzzz0zXXzOXXXX4X00zXzz0' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'idempotency-key: c_1308' \
  -d 'amount=500&destination%5Bamount%5D=500'
~~~

## License

CreditCardForm is available under the MIT license. See the LICENSE file for more info.
