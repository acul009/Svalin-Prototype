// Code generated by ent, DO NOT EDIT.

package user

import (
	"rahnit-rmm/ent/predicate"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.User {
	return predicate.User(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.User {
	return predicate.User(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.User {
	return predicate.User(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.User {
	return predicate.User(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.User {
	return predicate.User(sql.FieldLTE(FieldID, id))
}

// Username applies equality check predicate on the "username" field. It's identical to UsernameEQ.
func Username(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUsername, v))
}

// PasswordDoubleHashed applies equality check predicate on the "password_double_hashed" field. It's identical to PasswordDoubleHashedEQ.
func PasswordDoubleHashed(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPasswordDoubleHashed, v))
}

// Certificate applies equality check predicate on the "certificate" field. It's identical to CertificateEQ.
func Certificate(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldCertificate, v))
}

// PublicKey applies equality check predicate on the "public_key" field. It's identical to PublicKeyEQ.
func PublicKey(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPublicKey, v))
}

// EncryptedPrivateKey applies equality check predicate on the "encrypted_private_key" field. It's identical to EncryptedPrivateKeyEQ.
func EncryptedPrivateKey(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldEncryptedPrivateKey, v))
}

// TotpSecret applies equality check predicate on the "totp_secret" field. It's identical to TotpSecretEQ.
func TotpSecret(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldTotpSecret, v))
}

// UsernameEQ applies the EQ predicate on the "username" field.
func UsernameEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUsername, v))
}

// UsernameNEQ applies the NEQ predicate on the "username" field.
func UsernameNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldUsername, v))
}

// UsernameIn applies the In predicate on the "username" field.
func UsernameIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldUsername, vs...))
}

// UsernameNotIn applies the NotIn predicate on the "username" field.
func UsernameNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldUsername, vs...))
}

// UsernameGT applies the GT predicate on the "username" field.
func UsernameGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldUsername, v))
}

// UsernameGTE applies the GTE predicate on the "username" field.
func UsernameGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldUsername, v))
}

// UsernameLT applies the LT predicate on the "username" field.
func UsernameLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldUsername, v))
}

// UsernameLTE applies the LTE predicate on the "username" field.
func UsernameLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldUsername, v))
}

// UsernameContains applies the Contains predicate on the "username" field.
func UsernameContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldUsername, v))
}

// UsernameHasPrefix applies the HasPrefix predicate on the "username" field.
func UsernameHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldUsername, v))
}

// UsernameHasSuffix applies the HasSuffix predicate on the "username" field.
func UsernameHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldUsername, v))
}

// UsernameEqualFold applies the EqualFold predicate on the "username" field.
func UsernameEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldUsername, v))
}

// UsernameContainsFold applies the ContainsFold predicate on the "username" field.
func UsernameContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldUsername, v))
}

// PasswordDoubleHashedEQ applies the EQ predicate on the "password_double_hashed" field.
func PasswordDoubleHashedEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPasswordDoubleHashed, v))
}

// PasswordDoubleHashedNEQ applies the NEQ predicate on the "password_double_hashed" field.
func PasswordDoubleHashedNEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldPasswordDoubleHashed, v))
}

// PasswordDoubleHashedIn applies the In predicate on the "password_double_hashed" field.
func PasswordDoubleHashedIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldIn(FieldPasswordDoubleHashed, vs...))
}

// PasswordDoubleHashedNotIn applies the NotIn predicate on the "password_double_hashed" field.
func PasswordDoubleHashedNotIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldPasswordDoubleHashed, vs...))
}

// PasswordDoubleHashedGT applies the GT predicate on the "password_double_hashed" field.
func PasswordDoubleHashedGT(v []byte) predicate.User {
	return predicate.User(sql.FieldGT(FieldPasswordDoubleHashed, v))
}

// PasswordDoubleHashedGTE applies the GTE predicate on the "password_double_hashed" field.
func PasswordDoubleHashedGTE(v []byte) predicate.User {
	return predicate.User(sql.FieldGTE(FieldPasswordDoubleHashed, v))
}

// PasswordDoubleHashedLT applies the LT predicate on the "password_double_hashed" field.
func PasswordDoubleHashedLT(v []byte) predicate.User {
	return predicate.User(sql.FieldLT(FieldPasswordDoubleHashed, v))
}

// PasswordDoubleHashedLTE applies the LTE predicate on the "password_double_hashed" field.
func PasswordDoubleHashedLTE(v []byte) predicate.User {
	return predicate.User(sql.FieldLTE(FieldPasswordDoubleHashed, v))
}

// CertificateEQ applies the EQ predicate on the "certificate" field.
func CertificateEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldCertificate, v))
}

// CertificateNEQ applies the NEQ predicate on the "certificate" field.
func CertificateNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldCertificate, v))
}

// CertificateIn applies the In predicate on the "certificate" field.
func CertificateIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldCertificate, vs...))
}

// CertificateNotIn applies the NotIn predicate on the "certificate" field.
func CertificateNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldCertificate, vs...))
}

// CertificateGT applies the GT predicate on the "certificate" field.
func CertificateGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldCertificate, v))
}

// CertificateGTE applies the GTE predicate on the "certificate" field.
func CertificateGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldCertificate, v))
}

// CertificateLT applies the LT predicate on the "certificate" field.
func CertificateLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldCertificate, v))
}

// CertificateLTE applies the LTE predicate on the "certificate" field.
func CertificateLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldCertificate, v))
}

// CertificateContains applies the Contains predicate on the "certificate" field.
func CertificateContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldCertificate, v))
}

// CertificateHasPrefix applies the HasPrefix predicate on the "certificate" field.
func CertificateHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldCertificate, v))
}

// CertificateHasSuffix applies the HasSuffix predicate on the "certificate" field.
func CertificateHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldCertificate, v))
}

// CertificateEqualFold applies the EqualFold predicate on the "certificate" field.
func CertificateEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldCertificate, v))
}

// CertificateContainsFold applies the ContainsFold predicate on the "certificate" field.
func CertificateContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldCertificate, v))
}

// PublicKeyEQ applies the EQ predicate on the "public_key" field.
func PublicKeyEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPublicKey, v))
}

// PublicKeyNEQ applies the NEQ predicate on the "public_key" field.
func PublicKeyNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldPublicKey, v))
}

// PublicKeyIn applies the In predicate on the "public_key" field.
func PublicKeyIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldPublicKey, vs...))
}

// PublicKeyNotIn applies the NotIn predicate on the "public_key" field.
func PublicKeyNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldPublicKey, vs...))
}

// PublicKeyGT applies the GT predicate on the "public_key" field.
func PublicKeyGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldPublicKey, v))
}

// PublicKeyGTE applies the GTE predicate on the "public_key" field.
func PublicKeyGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldPublicKey, v))
}

// PublicKeyLT applies the LT predicate on the "public_key" field.
func PublicKeyLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldPublicKey, v))
}

// PublicKeyLTE applies the LTE predicate on the "public_key" field.
func PublicKeyLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldPublicKey, v))
}

// PublicKeyContains applies the Contains predicate on the "public_key" field.
func PublicKeyContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldPublicKey, v))
}

// PublicKeyHasPrefix applies the HasPrefix predicate on the "public_key" field.
func PublicKeyHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldPublicKey, v))
}

// PublicKeyHasSuffix applies the HasSuffix predicate on the "public_key" field.
func PublicKeyHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldPublicKey, v))
}

// PublicKeyEqualFold applies the EqualFold predicate on the "public_key" field.
func PublicKeyEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldPublicKey, v))
}

// PublicKeyContainsFold applies the ContainsFold predicate on the "public_key" field.
func PublicKeyContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldPublicKey, v))
}

// EncryptedPrivateKeyEQ applies the EQ predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldEncryptedPrivateKey, v))
}

// EncryptedPrivateKeyNEQ applies the NEQ predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyNEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldEncryptedPrivateKey, v))
}

// EncryptedPrivateKeyIn applies the In predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldIn(FieldEncryptedPrivateKey, vs...))
}

// EncryptedPrivateKeyNotIn applies the NotIn predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyNotIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldEncryptedPrivateKey, vs...))
}

// EncryptedPrivateKeyGT applies the GT predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyGT(v []byte) predicate.User {
	return predicate.User(sql.FieldGT(FieldEncryptedPrivateKey, v))
}

// EncryptedPrivateKeyGTE applies the GTE predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyGTE(v []byte) predicate.User {
	return predicate.User(sql.FieldGTE(FieldEncryptedPrivateKey, v))
}

// EncryptedPrivateKeyLT applies the LT predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyLT(v []byte) predicate.User {
	return predicate.User(sql.FieldLT(FieldEncryptedPrivateKey, v))
}

// EncryptedPrivateKeyLTE applies the LTE predicate on the "encrypted_private_key" field.
func EncryptedPrivateKeyLTE(v []byte) predicate.User {
	return predicate.User(sql.FieldLTE(FieldEncryptedPrivateKey, v))
}

// TotpSecretEQ applies the EQ predicate on the "totp_secret" field.
func TotpSecretEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldTotpSecret, v))
}

// TotpSecretNEQ applies the NEQ predicate on the "totp_secret" field.
func TotpSecretNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldTotpSecret, v))
}

// TotpSecretIn applies the In predicate on the "totp_secret" field.
func TotpSecretIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldTotpSecret, vs...))
}

// TotpSecretNotIn applies the NotIn predicate on the "totp_secret" field.
func TotpSecretNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldTotpSecret, vs...))
}

// TotpSecretGT applies the GT predicate on the "totp_secret" field.
func TotpSecretGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldTotpSecret, v))
}

// TotpSecretGTE applies the GTE predicate on the "totp_secret" field.
func TotpSecretGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldTotpSecret, v))
}

// TotpSecretLT applies the LT predicate on the "totp_secret" field.
func TotpSecretLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldTotpSecret, v))
}

// TotpSecretLTE applies the LTE predicate on the "totp_secret" field.
func TotpSecretLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldTotpSecret, v))
}

// TotpSecretContains applies the Contains predicate on the "totp_secret" field.
func TotpSecretContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldTotpSecret, v))
}

// TotpSecretHasPrefix applies the HasPrefix predicate on the "totp_secret" field.
func TotpSecretHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldTotpSecret, v))
}

// TotpSecretHasSuffix applies the HasSuffix predicate on the "totp_secret" field.
func TotpSecretHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldTotpSecret, v))
}

// TotpSecretEqualFold applies the EqualFold predicate on the "totp_secret" field.
func TotpSecretEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldTotpSecret, v))
}

// TotpSecretContainsFold applies the ContainsFold predicate on the "totp_secret" field.
func TotpSecretContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldTotpSecret, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.User) predicate.User {
	return predicate.User(sql.NotPredicates(p))
}
