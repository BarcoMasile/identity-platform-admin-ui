import React, { FC } from "react";
import {
  ActionButton,
  Button,
  Col,
  Row,
  useNotify,
} from "@canonical/react-components";
import { useFormik } from "formik";
import * as Yup from "yup";
import { queryKeys } from "util/queryKeys";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import ClientForm, { ClientFormTypes } from "pages/clients/ClientForm";
import { fetchClient, updateClient } from "api/client";
import usePanelParams from "util/usePanelParams";
import SidePanel from "components/SidePanel";

const ClientEdit: FC = () => {
  const navigate = useNavigate();
  const notify = useNotify();
  const queryClient = useQueryClient();
  const panelParams = usePanelParams();
  const clientId = panelParams.id;

  if (!clientId) {
    return;
  }

  const { data: client } = useQuery({
    queryKey: [queryKeys.clients, clientId],
    queryFn: () => fetchClient(clientId),
  });

  const ClientEditSchema = Yup.object().shape({
    client_name: Yup.string().required("This field is required"),
  });

  const formik = useFormik<ClientFormTypes>({
    initialValues: {
      client_uri: client?.client_uri,
      client_name: client?.client_name,
      grant_types: client?.grant_types,
      response_types: client?.response_types,
      scope: client?.scope,
      redirect_uris: client?.redirect_uris,
      request_object_signing_alg: client?.request_object_signing_alg,
    },
    enableReinitialize: true,
    validationSchema: ClientEditSchema,
    onSubmit: (values) => {
      updateClient(client?.client_id ?? "", JSON.stringify(values))
        .then(() => {
          void queryClient.invalidateQueries({
            queryKey: [queryKeys.clients],
          });
          navigate(`/client`, notify.queue(notify.success("Client updated")));
        })
        .catch((e) => {
          formik.setSubmitting(false);
          notify.failure("Client update failed", e);
        });
    },
  });

  const submitForm = () => {
    void formik.submitForm();
  };

  return (
    <SidePanel hasError={false} loading={false} className="p-panel">
      <SidePanel.Header>
        <SidePanel.HeaderTitle>Edit client</SidePanel.HeaderTitle>
      </SidePanel.Header>
      <SidePanel.Content className="u-no-padding">
        <Row>
          <ClientForm formik={formik} />
        </Row>
      </SidePanel.Content>
      <SidePanel.Footer className="u-align--right">
        <Row>
          <Col size={12}>
            <Button
              appearance="base"
              className="u-no-margin--bottom u-sv2"
              onClick={() => navigate(`/client`)}
            >
              Cancel
            </Button>
            <ActionButton
              appearance="positive"
              className="u-no-margin--bottom"
              loading={formik.isSubmitting}
              disabled={!formik.isValid}
              onClick={submitForm}
            >
              Update
            </ActionButton>
          </Col>
        </Row>
      </SidePanel.Footer>
    </SidePanel>
  );
};

export default ClientEdit;
