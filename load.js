import { Client,StatusOK } from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new Client();
client.load(['.'], 'def.proto');





export default () => {
    client.connect('localhost:50001',{ plaintext: true });

    for( let i =0; i<= 100; i++){
      const data = { valueToIncrement: getRndInteger(1,10) };

      const response = client.invoke('SerializedService/ParentUnaryCall', data);
    
      check(response, {
        'status is OK': (r) => r && r.status === StatusOK,
      });
    
    }
    
    sleep(1);
  };
  


function getRndInteger(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}