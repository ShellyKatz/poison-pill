package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type TestDataItem struct {
	fieldName string
	fieldValue string
	fieldMinValue string
	errorMessage string
}

var _ = Describe("PoisonPillConfig Validation", func() {

	Describe("validate time fields of PoisonPillConfig CR", func() {
		shortTimeTestItems := []TestDataItem{
			{"PeerApiServerTimeout", "1.2ms", MinDurPeerApiServerTimeout, ErrPeerApiServerTimeout},
			{"ApiServerTimeout", "1.2ms", MinDurApiServerTimeout, ErrApiServerTimeout},
			{"PeerDialTimeout", "1.2ms", MinDurPeerDialTimeout, ErrPeerDialTimeout},
			{"PeerRequestTimeout", "1.2ms", MinDurPeerRequestTimeout, ErrPeerRequestTimeout},
			{"ApiCheckInterval", "1.2ms", MinDurApiCheckInterval, ErrApiCheckInterval},
			{"PeerUpdateInterval", "1.2ms", MinDurPeerUpdateInterval, ErrPeerUpdateInterval},
		}

		for _, item := range shortTimeTestItems {
			item := item
			text := "for " + item.fieldName + " value shorter than " + item.fieldMinValue
			Context(text, func() {
				It("should be rejected", func() {
					ppc := createPoisonPillConfigCR(item.fieldName, item.fieldValue)
					//todo- question: is there a better option to change a field value when the field is a parameter
					//fmt.Printf("%+v,\n %s\n",ppc, item.fieldName)
					err := ppc.ValidateTimes()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()) .To(ContainSubstring(item.errorMessage))
				})
			})
		}

		negativeTimeTestItems := []TestDataItem{
			{"PeerApiServerTimeout", "-1ms", MinDurPeerApiServerTimeout, ErrPeerApiServerTimeout},
			{"ApiServerTimeout", "-1ms", MinDurApiServerTimeout, ErrApiServerTimeout},
			{"PeerDialTimeout", "-1ms", MinDurPeerDialTimeout, ErrPeerDialTimeout},
			{"PeerRequestTimeout", "-1ms", MinDurPeerRequestTimeout, ErrPeerRequestTimeout},
			{"ApiCheckInterval", "-1ms", MinDurApiCheckInterval, ErrApiCheckInterval},
			{"PeerUpdateInterval", "-1ms", MinDurPeerUpdateInterval, ErrPeerUpdateInterval},
		}

		for _, item := range negativeTimeTestItems {
			item := item
			text := "for " + item.fieldName + " with negative value" + item.fieldMinValue
			Context(text, func() {
				It("should be rejected", func() {
					ppc := createPoisonPillConfigCR(item.fieldName, item.fieldValue)
					err := ppc.ValidateTimes()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()) .To(ContainSubstring(item.errorMessage))
				})
			})
		}

		Context("for valid CR", func() {
			It("should not be rejected", func() {
				ppc := createDefaultPoisonPillConfigCR()
				err := ppc.ValidateTimes()
				Expect(err).NotTo(HaveOccurred())

			})
		})

	})

	Describe("parsing time duration string to milliseconds", func() {
		Context("for the string 10ms", func(){
			It("should return 10", func() {
					Expect(toMS("10ms")).To(Equal(int64(10)))
			})
		})

		Context("for the string 10s", func(){
			It("should return 10000", func(){
				Expect(toMS("10s")).To(Equal(int64(10000)))
			})
		})

	})
})

func createDefaultPoisonPillConfigCR() *PoisonPillConfig {
	ppc := &PoisonPillConfig{}
	ppc.Name = "test"
	ppc.Namespace = "default"

	//default values for time fields
	ppc.Spec.PeerApiServerTimeout = &metav1.Duration{Duration: 5*time.Second}
	ppc.Spec.ApiServerTimeout = &metav1.Duration{Duration: 5*time.Second}
	ppc.Spec.PeerDialTimeout = &metav1.Duration{Duration: 5*time.Second}
	ppc.Spec.PeerRequestTimeout = &metav1.Duration{Duration: 5*time.Second}
	ppc.Spec.ApiCheckInterval = &metav1.Duration{Duration: 15*time.Second}
	ppc.Spec.PeerUpdateInterval = &metav1.Duration{Duration: 15*time.Minute}

	return ppc
}

func createPoisonPillConfigCR(fieldName string, value string) *PoisonPillConfig {
	ppc := createDefaultPoisonPillConfigCR()

	//set the field tested
	setFieldValue(ppc, fieldName, value)

	return ppc
}



func setFieldValue(ppc *PoisonPillConfig, fieldName string, value string) {
	d, err := time.ParseDuration(value)
	if err != nil {
		//todo return error!
	}
	timeValue := &metav1.Duration{Duration: d}
	switch fieldName {
	case "PeerApiServerTimeout":
		ppc.Spec.PeerApiServerTimeout = timeValue
	case "ApiServerTimeout":
		ppc.Spec.ApiServerTimeout = timeValue
	case "PeerDialTimeout":
		ppc.Spec.PeerDialTimeout = timeValue
	case "PeerRequestTimeout":
		ppc.Spec.PeerRequestTimeout = timeValue
	case "ApiCheckInterval":
		ppc.Spec.ApiCheckInterval = timeValue
	case "PeerUpdateInterval":
		ppc.Spec.PeerUpdateInterval = timeValue
	}

}

//todo - check on smaller times, float times and negative? and on correct times! (bigger or equal)