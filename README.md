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

``` sh
brew install go --cross-compile-common
```

Or if you have Go installed you can build it:

``` sh
go get 

go build
```

Run it:

``` sh
mockstripe
```

Then from another terminal:

``` sh
curl -i http://localhost:8080/v1/version -H "Authorization: Bearer sk_test_123"
```

By default, stripe-mock runs on port 8080, but is configurable with the
`-port` option. (TODO)

## Development
TBD

## Supported Operations

- [x] Auth
- [x] Capture
- [x] Refund

## License

CreditCardForm is available under the MIT license. See the LICENSE file for more info.