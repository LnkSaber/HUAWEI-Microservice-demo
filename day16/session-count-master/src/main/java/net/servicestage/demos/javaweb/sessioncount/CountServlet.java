package net.servicestage.demos.javaweb.sessioncount;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.logging.Logger;

public class CountServlet extends HttpServlet {

    private final static String SESSION_ID_INTEGER = "int-count";
    private static Logger logger = Logger.getLogger(CountServlet.class.getName());

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp)
            throws ServletException, IOException {

        Integer count = (Integer) req.getSession(true).getAttribute(SESSION_ID_INTEGER);
        logger.info( "Count object in session:" + count );
        if(count == null) {
            count = 0;
        }
        count++;
        req.getSession().setAttribute( SESSION_ID_INTEGER, count );

        resp.setStatus( 200 );
        resp.setHeader( "Count", count.toString() );
        resp.getWriter().append( "Hello ServiceStage, now the count is [" ).append( count.toString() ).append( "]. <BR> Please refresh this page." ).flush();
    }
}