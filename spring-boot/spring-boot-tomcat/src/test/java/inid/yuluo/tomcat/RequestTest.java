package inid.yuluo.tomcat;

import java.io.IOException;
import java.net.HttpURLConnection;
import java.net.URL;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */
public class RequestTest {

	public static void main(String[] args) {

		for (int i = 0; i < 10; i++) {
			new Thread(new RequestTask()).start();
		}
	}

	private static class RequestTask implements Runnable {
		@Override
		public void run() {
			try {
				URL url = new URL("http://localhost:8080/");
				HttpURLConnection connection = (HttpURLConnection) url.openConnection();
				connection.setRequestMethod("GET");
				int responseCode = connection.getResponseCode();
				System.out.println("Response Code: " + responseCode);
			}
			catch (IOException e) {
				e.printStackTrace();
			}
		}
	}

}

