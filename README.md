# gomendan assistant

This is a Probe Of Concept about showing the face emotion of a person using the webcam of the computer.

## Getting Started

Clone the repository:

```bash
git clone https://github.com/ervitis/gomendan-assistant
```

This application uses [Google Vision API](https://cloud.google.com/vision?hl=en), so you need to create a [Service Account](https://cloud.google.com/iam/docs/service-account-overview) and download the JSON file which includes the key credentials to connect to the API.

> :warning: **Be aware of not sharing that file or include it in the repository.**

Then export an environment variable:

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/absolute/path/to/the/credentials/json/file
```

Finally, run the application:

```bash
make run
```

### Prerequisites

[gocv](https://github.com/hybridgroup/gocv/tree/release) uses opencv internally, so you need to install all the dependencies from [here](https://github.com/hybridgroup/gocv/tree/release#how-to-install).

Then, you execute the command

```bash
make run
```

And it will download all the dependencies

## Built With

* [gocv](ttps://github.com/hybridgroup/gocv) - GoCV
* [Google Vision API](https://pkg.go.dev/cloud.google.com/go/vision/v2/apiv1) - For face emotion analysis

## Authors

* **ervitis** - *Initial work* - [ervitis](https://github.com/ervitis)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Ron Evans for his great work developing this library: https://www.youtube.com/watch?v=mId8cX4h_Ms
